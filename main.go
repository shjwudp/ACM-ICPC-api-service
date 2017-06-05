package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	gin_gzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

// Env is all handlers common environment
type Env struct {
	db         model.Datastore
	printQueue chan int64
}

func (env *Env) initAdmin(account, password string) {
	admin := model.User{
		Account:  account,
		Password: password,
		Role:     "admin",
	}
	env.db.SaveUser(admin)
}

func (env *Env) sendToPrint(printerName string) {
	errorHandler := func(id int64, err error) {
		log.Println(err)
		env.printQueue <- id
	}

	var id int64
	ok := true
	for ok {
		if id, ok = <-env.printQueue; ok {
			p, err := env.db.GetPrintCode(id)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetPrintCode failed with %s", err))
				continue
			}
			u, err := env.db.GetUserAccount(p.Account)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetUserAccount failed with %s", err))
				continue
			}

			cmd := exec.Command("lp",
				"-d", printerName,
				"-t", u.DisplayName+"-"+u.SeatID,
				"-o", "prettyprint",
				"-o", "Page-left=36",
				"-o", "Page-right=36",
				"-o", "Page-top=36",
				"-o", "Page-bottom=36")
			cmd.Stdin = strings.NewReader(p.Code)
			err = cmd.Run()
			if err != nil {
				errorHandler(id, fmt.Errorf("Run cmd lp failed with %s", err))
				continue
			}
			p.IsDone = true
			err = env.db.UpdatePrintCode(*p)
			if err != nil {
				errorHandler(id, fmt.Errorf("Run cmd lp failed with %s", err))
			}
		}
	}
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
	Config, err = model.ConfigurationLoad(confPath)
	if err != nil {
		log.Fatal("Load conf failed with", err)
	}

	db, err := model.OpenDB(Config.Storage.Dirver, Config.Storage.Config)
	if err != nil {
		log.Fatal("Open db failed with", err)
	}

	// init env
	env := &Env{
		db:         db,
		printQueue: make(chan int64, Config.Printer.QueueSize),
	}
	// init admin
	env.initAdmin(Config.Server.Admin.Account, Config.Server.Admin.Password)
	// start a goruntine to update ContestStanding
	go func() {
		for {
			env.updateContestStanding("./results.xml")
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	// start a group of goruntine to deal with print task
	for _, name := range Config.Printer.PinterNameList {
		go env.sendToPrint(name)
	}

	// gin.SetMode(gin.TestMode)
	gin.SetMode(gin.ReleaseMode)
	router := GetMainEngine(env, Config.Server.JWTSecret)
	s := &http.Server{
		Addr:    Config.Server.Port,
		Handler: router,
	}
	s.ListenAndServe()
}

func (env *Env) getContestStanding(c *gin.Context) {
	kv, err := env.db.GetKV("ContestStanding")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
	}
	c.String(200, string(kv.Value))
}

func (env *Env) listBallonStatus(c *gin.Context) {
	// update BallonStatus first
	kv, err := env.db.GetKV("ContestStanding")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
	}

	var cs model.ContestStanding
	err = json.Unmarshal(kv.Value, &cs)
	if err != nil {
		var errMsg = fmt.Sprint("Json Unmarshal failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
		return
	}

	var results []map[string]interface{}
	for _, t := range cs.TeamStandings {
		for _, p := range t.ProblemSummaryInfos {
			if !p.IsSolved {
				continue
			}
			team, err := env.db.GetUserTeamKey(t.TeamKey)
			if err != nil {
				log.Printf("GetUserTeamKey %s failed with %s\n", t.TeamKey, err)
				continue
			}
			bs, err := env.db.GetBallonStatus(t.TeamKey, p.Index)
			if err == sql.ErrNoRows {
				bs = &model.BallonStatus{
					TeamKey:      t.TeamKey,
					ProblemIndex: p.Index,
					IsMarked:     false,
				}
			} else if err != nil {
				errMsg := fmt.Sprint("GetBallonStatus failed with", err)
				c.JSON(500, gin.H{"message": errMsg})
				return
			}
			results = append(results, map[string]interface{}{
				"TeamKey":       t.TeamKey,
				"ProblemIndex":  p.Index,
				"SolutionTime":  p.SolutionTime,
				"SeatID":        team.SeatID,
				"IsSolved":      p.IsSolved,
				"IsFirstSolved": p.IsFirstSolved,
				"IsMarked":      bs.IsMarked,
			})
		}
	}
	data, _ := json.Marshal(results)
	c.String(http.StatusOK, string(data))
}

func (env *Env) patchBallonStatus(c *gin.Context) {
	var req struct {
		TeamKey      string `json:"TeamKey" binding:"required"`
		ProblemIndex int    `json:"ProblemIndex" binding:"required"`
		Action       string `json:"action" binding:"required"`
	}
	err := c.BindJSON(&req)
	if err == nil {
		log.Println(req)
		if req.Action != "mark" {
			errMsg := fmt.Sprintf("No such action=%s", req.Action)
			c.JSON(400, gin.H{"message": errMsg})
			return
		}
		err := env.db.SaveBallonStatus(model.BallonStatus{
			TeamKey:      req.TeamKey,
			ProblemIndex: req.ProblemIndex,
			IsMarked:     true,
		})
		if err != nil {
			errMsg := fmt.Sprintf("No such action=%s", req.Action)
			c.JSON(400, gin.H{"message": errMsg})
			return
		}
		c.JSON(200, gin.H{"message": "OK"})
		return
	}
	c.JSON(400, gin.H{"message": fmt.Sprint("BindJSON failed with", err)})
}

func (env *Env) postPrinter(c *gin.Context) {
	var requestJSON struct {
		PrintContent string `json:"PrintContent" binding:"required"`
	}
	err := c.BindJSON(&requestJSON)
	if err != nil {
		errMsg := fmt.Sprint("BindJSON failed with", err)
		c.JSON(400, gin.H{"message": errMsg})
		return
	}

	account, _ := c.Get("account")
	p, err := env.db.SavePrintCode(account.(string), requestJSON.PrintContent)
	if err != nil {
		errMsg := fmt.Sprint("SavePrintCode failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
		return
	}
	env.printQueue <- p.ID
	c.JSON(200, gin.H{"message": "OK", "queue_size": len(env.printQueue)})
}

func (env *Env) getParticipant(c *gin.Context) {
	users, err := env.db.ListUser()
	if err != nil {
		errMsg := fmt.Sprint("List User failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
		return
	}
	data, err := json.Marshal(users)
	if err != nil {
		errMsg := fmt.Sprint("List User failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
		return
	}
	c.String(http.StatusOK, string(data))
}

func (env *Env) postParticipant(c *gin.Context) {
	var err error
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		errMsg := fmt.Sprint("Get uploadFile failed with", err)
		c.JSON(400, gin.H{"message": errMsg})
		return
	}
	filename := header.Filename
	log.Println("filename : ", filename)

	titleMap := make(map[string]int)
	r := csv.NewReader(file)
	r.Comma = '\t'
	r.Comment = '#'
	lineNo := 1
	for {
		A, err := r.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
			c.JSON(400, gin.H{"message": errMsg})
			return
		}
		if lineNo == 1 {
			for colNo, title := range A {
				titleMap[title] = colNo
			}
			log.Println(titleMap)
			if _, ok := titleMap["account"]; !ok {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(400, gin.H{"message": errMsg})
				return
			}
			if _, ok := titleMap["site"]; !ok {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(400, gin.H{"message": errMsg})
				return
			}
		} else {
			if len(A) != len(titleMap) {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(400, gin.H{"message": errMsg})
				return
			}
			u := model.User{
				Account: A[titleMap["account"]],
				TeamKey: strings.ToUpper(A[titleMap["site"]] + A[titleMap["account"]]),
			}
			if i, ok := titleMap["password"]; ok {
				u.Password = A[i]
			}
			if i, ok := titleMap["displayname"]; ok {
				u.DisplayName = A[i]
			}
			if i, ok := titleMap["nickname"]; ok {
				u.NickName = A[i]
			}
			if i, ok := titleMap["school"]; ok {
				u.School = A[i]
			}
			if i, ok := titleMap["isstar"]; ok {
				u.IsStar = A[i]
			}
			if i, ok := titleMap["isgirl"]; ok {
				u.IsGirl = A[i]
			}
			if i, ok := titleMap["seatid"]; ok {
				u.SeatID = A[i]
			}
			if i, ok := titleMap["coach"]; ok {
				u.Coach = A[i]
			}
			if i, ok := titleMap["player1"]; ok {
				u.Player1 = A[i]
			}
			if i, ok := titleMap["player2"]; ok {
				u.Player2 = A[i]
			}
			if i, ok := titleMap["player3"]; ok {
				u.Player3 = A[i]
			}
			err := env.db.SaveUser(u)
			if err != nil {
				errMsg := fmt.Sprint("Saved user failed with", err)
				c.JSON(400, gin.H{"message": errMsg})
				return
			}
		}
		lineNo++
	}
	c.JSON(200, gin.H{"message": "OK"})
}

func (env *Env) genPostAuthenticate(jwtSecret string) func(*gin.Context) {
	return func(c *gin.Context) {
		var requestJSON struct {
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		log.Println(c.Request)
		if c.BindJSON(&requestJSON) == nil {
			log.Println("requestJSON :", requestJSON)
			user, err := env.db.GetUserAccount(requestJSON.Account)
			if err != nil {
				errMsg := fmt.Sprint("Get User failed with", err)
				c.JSON(500, gin.H{"message": errMsg})
				return
			}
			if user.Password == requestJSON.Password {
				log.Println("user :", user)
				// Create the token
				token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
				// Set some claims
				token.Claims = jwt_lib.MapClaims{
					"Account": user.Account,
					"Role":    user.Role,
					"exp":     time.Now().Add(time.Hour * 12).Unix(),
				}
				// Sign and get the complete encoded token as a string
				tokenString, err := token.SignedString([]byte(jwtSecret))
				if err != nil {
					c.JSON(500, gin.H{"message": "Could not generate token"})
				}
				c.JSON(200, gin.H{"token": tokenString})
				return
			}
			c.JSON(401, gin.H{"message": "Wrong account or password"})
			return
		}
		c.JSON(401, gin.H{"message": "Invalid Request"})
	}
}

// FakeJWTAuthMiddleware : for performance test
func FakeJWTAuthMiddleware(c *gin.Context) {
	c.Set("account", "admin")
	c.Set("role", "admin")
	c.Next()
}

// JWTAuthMiddleware : JWT Authorization Verification
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt_request.ParseFromRequest(
			c.Request,
			jwt_request.OAuth2Extractor,
			func(token *jwt_lib.Token) (interface{}, error) {
				b := ([]byte(secret))
				return b, nil
			},
		)

		c.Set("token", token)
		account, ok := token.Claims.(jwt_lib.MapClaims)["Account"]
		if !ok {
			c.AbortWithError(500, errors.New("No Account in token"))
			return
		}
		c.Set("account", account)
		role, ok := token.Claims.(jwt_lib.MapClaims)["Role"]
		if !ok {
			c.AbortWithError(500, errors.New("No Role in token"))
			return
		}
		c.Set("role", role)

		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}

// CORSMiddleware : CORS Middleware
func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}
	c.Next()
}

var level2Role = map[string]bool{
	"normal": true,
}

var level1Role = map[string]bool{
	"volunteer": true,
}

var level0Role = map[string]bool{
	"admin": true,
}

func level1PermissionMiddleware(c *gin.Context) {
	role, _ := c.Get("role")
	if _, ok := level0Role[role.(string)]; ok {
		c.Next()
		return
	}
	if _, ok := level1Role[role.(string)]; ok {
		c.Next()
		return
	}
	c.AbortWithError(403, errors.New("No Permission"))
}

func level0PermissionMiddleware(c *gin.Context) {
	role, _ := c.Get("role")
	if _, ok := level0Role[role.(string)]; ok {
		c.Next()
		return
	}
	c.AbortWithError(403, errors.New("No Permission"))
}

// GetMainEngine : Main Engine
func GetMainEngine(env *Env, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware)
	router.Use(gin_gzip.Gzip(gin_gzip.DefaultCompression))

	api := router.Group("/api")
	{
		postAuthenticate := env.genPostAuthenticate(jwtSecret)
		api.POST("/authenticate", postAuthenticate)

		authorized := api.Group("/authorized")
		authorized.Use(JWTAuthMiddleware(jwtSecret))
		// authorized.Use(FakeJWTAuthMiddleware)
		{
			authorized.GET("/contest-standing", env.getContestStanding)
			authorized.POST("/printer", env.postPrinter)
			level1 := authorized.Group("/", level1PermissionMiddleware)
			{
				level1.GET("/ballon-status", env.listBallonStatus)
				level1.PATCH("/ballon-status", env.patchBallonStatus)
				level0 := authorized.Group("/", level0PermissionMiddleware)
				{
					level0.GET("/participant", env.getParticipant)
					level0.POST("/participant", env.postParticipant)
				}
			}
			// authorized.GET("/ballon-status", env.listBallonStatus)
			// authorized.PATCH("/ballon-status", env.patchBallonStatus)
			// authorized.GET("/participant", env.getParticipant)
			// authorized.POST("/participant", env.postParticipant)
		}
	}
	return router
}

func (env *Env) updateContestStanding(resultsXMLPath string) {
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
	err = env.db.SaveKV(kv)
	if err != nil {
		log.Println("SaveKV failed with", err)
		return
	}
}
