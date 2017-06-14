package main

import (
	"encoding/json"
	"flag"
	gin_gzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"github.com/shjwudp/ACM-ICPC-api-service/server"
	"github.com/shjwudp/ACM-ICPC-api-service/server/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

var usageStr = `
Usage: ACM-ICPC-api-service --config <file>		Configuration file path
`

// usage will print out the flag options for the server.
func usage() {
	log.Println(usageStr)
	os.Exit(0)
}

func initWithConf(conf model.Configuration) (*server.Env, error) {
	db, err := model.OpenDB(conf.Storage)
	if err != nil {
		return nil, err
	}

	// init env
	env := server.NewEnv(db,
		conf.Printer.QueueSize,
		conf.Server.JWTSecret)

	initAdmin(db,
		conf.Server.Admin.Account,
		conf.Server.Admin.Password)

	// start a goruntine to update ContestStanding
	go func() {
		for {
			updateContestStanding(db, conf.ResultsXMLPath)
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	// start a group of goruntine to deal with print task
	for _, name := range conf.Printer.PinterNameList {
		go env.SendPrinter(name)
	}
	b, err := json.Marshal(conf.ContestInfo)
	if err != nil {
		return nil, err
	}
	err = db.SaveKV(model.KV{Key: "ContestInfo", Value: b})
	if err != nil {
		return nil, err
	}
	return env, nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var confPath string
	flag.StringVar(&confPath, "config", "", "Configuration file path.")
	flag.Parse()
	flag.Usage = usage
	flag.Parse()
	if confPath == "" {
		usage()
	}

	var err error
	conf, err := model.ConfigurationLoad(confPath)
	if err != nil {
		log.Fatal("Load conf failed with", err)
	}
	env, err := initWithConf(*conf)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Server.IsTestMode {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := GetMainEngine(env, conf.Server.JWTSecret, conf.Server.IsTestMode)
	s := &http.Server{
		Addr:    conf.Server.Addr,
		Handler: router,
	}
	s.ListenAndServe()
}

// GetMainEngine : Main Engine
func GetMainEngine(env *server.Env, JWTSecret string, NeedAuth bool) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Options)
	router.Use(gin_gzip.Gzip(gin_gzip.DefaultCompression))

	api := router.Group("/api")
	{
		api.POST("/authenticate", env.Authenticate)

		authorized := api.Group("/authorized")
		if NeedAuth {
			authorized.Use(middleware.JWTAuthMiddleware(JWTSecret))
		} else {
			authorized.Use(middleware.FakeJWTAuthMiddleware)
		}
		{
			authorized.GET("/contest-standing", env.GetContestStanding)
			authorized.POST("/printer", env.PostPrinter)
			level1 := authorized.Group("/", middleware.Level1PermissionMiddleware)
			{
				level1.GET("/ballon-status", env.ListBallonStatus)
				level1.PATCH("/ballon-status", env.PatchBallonStatus)
				level0 := authorized.Group("/", middleware.Level0PermissionMiddleware)
				{
					level0.GET("/participant", env.AllUser)
					level0.POST("/participant", env.PostUserList)
					level0.GET("/contest-info", env.GetContestInfo)
					level0.POST("/contest-info", env.SaveContestInfo)
				}
			}
		}
	}
	return router
}

func updateContestStanding(db model.Datastore, resultsXMLPath string) {
	cs, err := model.ParseResultXML(resultsXMLPath)
	if err != nil {
		log.Println("ParseResultXML failed with", err)
		return
	}
	allUser, err := db.AllUser()
	if err != nil {
		log.Println("ParseResultXML failed with", err)
		return
	}
	userMap := map[string]model.User{}
	for _, user := range allUser {
		userMap[user.TeamKey] = user
	}

	for _, t := range cs.TeamStandings {
		if u, has := userMap[t.TeamKey]; has {
			t.NickName = u.NickName
			t.School = u.School
			t.IsStar = u.IsStar
			t.IsGirl = u.IsGirl
			t.Coach = u.Coach
			t.Player1 = u.Player1
			t.Player2 = u.Player2
			t.Player3 = u.Player3
			t.SeatID = u.SeatID
		}
		for _, p := range t.ProblemSummaryInfos {
			db.SaveContestEvent(model.ContestEvent{
				TimeStamp:    cs.StandingsHeader.CurrentTimeStamp,
				TeamKey:      t.TeamKey,
				ProblemIndex: p.Index,
				Attempts:     p.Attempts,
				IsSolved:     p.IsSolved,
				Points:       p.Points,
			})
		}
	}

	b, err := json.Marshal(cs)
	if err != nil {
		log.Println("json.Marshal failed with", err)
		return
	}
	kv := model.KV{
		Key:   "ContestStanding",
		Value: b,
	}
	log.Println("Update ContestStanding")
	err = db.SaveKV(kv)
	if err != nil {
		log.Println("SaveKV failed with", err)
		return
	}
}

func initAdmin(db model.Datastore, account, password string) error {
	admin := model.User{
		Account:  account,
		Password: password,
		Role:     "admin",
	}
	return db.SaveUser(admin)
}
