package main

import (
	"encoding/json"
	"flag"
	gin_gzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

// Config is configuration struct
var Config *model.Configuration

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
	Config, err = model.ConfigurationLoad(confPath)
	if err != nil {
		log.Fatal("Load conf failed with", err)
	}

	db, err := model.OpenDB(Config.Storage.Dirver, Config.Storage.Config)
	if err != nil {
		log.Fatal("Open db failed with", err)
	}

	// init env
	env := server.NewEnv(
		db,
		Config.Printer.QueueSize,
		Config.Server.JWTSecret)

	initAdmin(db,
		Config.Server.Admin.Account,
		Config.Server.Admin.Password)

	// start a goruntine to update ContestStanding
	go func() {
		for {
			updateContestStanding(db, Config.ResultsXMLPath)
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	// start a group of goruntine to deal with print task
	for _, name := range Config.Printer.PinterNameList {
		go env.SendPrinter(name)
	}

	gin.SetMode(gin.TestMode)
	// gin.SetMode(gin.ReleaseMode)
	router := GetMainEngine(env, Config.Server.JWTSecret)
	s := &http.Server{
		Addr:    Config.Server.Port,
		Handler: router,
	}
	s.ListenAndServe()
}

// GetMainEngine : Main Engine
func GetMainEngine(env *server.Env, JWTSecret string) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Options)
	router.Use(gin_gzip.Gzip(gin_gzip.DefaultCompression))

	api := router.Group("/api")
	{
		api.POST("/authenticate", env.Authenticate)

		authorized := api.Group("/authorized")
		authorized.Use(middleware.JWTAuthMiddleware(JWTSecret))
		// authorized.Use(FakeJWTAuthMiddleware)
		{
			authorized.GET("/contest-standing", env.GetContestStanding)
			authorized.POST("/printer", env.PostPrinter)
			level1 := authorized.Group("/", middleware.Level1PermissionMiddleware)
			{
				level1.GET("/ballon-status", env.ListBallonStatus)
				level1.PATCH("/ballon-status", env.PatchBallonStatus)
				level0 := authorized.Group("/", middleware.Level0PermissionMiddleware)
				{
					level0.GET("/participant", env.ListUser)
					level0.POST("/participant", env.PostUserList)
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
	b, _ := json.Marshal(cs)
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

func initAdmin(db model.Datastore, account, password string) {
	admin := model.User{
		Account:  account,
		Password: password,
		Role:     "admin",
	}
	db.SaveUser(admin)
}
