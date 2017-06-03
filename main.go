package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/shjwudp/ACM-ICPC-api-service/model"

	jwt_lib "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gin-gonic/gin.v1"
	"strings"
)

var usageStr = `
Usage: ACM-ICPC-api-service --config <file>		Configuration file path
`
var dbEngine *xorm.Engine

// usage will print out the flag options for the server.
func usage() {
	log.Println(usageStr)
	os.Exit(0)
}

// Config is configuration struct
var Config = struct {
	Server struct {
		JWTSecret string
		Port      string
		Admin     struct {
			Account  string
			Password string
		}
	}
	// use sqlite3
	Storage struct {
		Dirver string
		Config string
	}
	ResultsXMLPath string
}{}

func loadConfJSON(confPath string) error {
	configFile, err := os.Open(confPath)

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&Config)

	return err
}

func main() {
	var confPath string
	flag.StringVar(&confPath, "config", "", "Configuration file path.")
	flag.Parse()

	flag.Usage = usage
	flag.Parse()
	if confPath == "" {
		usage()
	}

	var err = loadConfJSON(confPath)
	if err != nil {
		log.Fatal("Load conf failed")
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err = setupDB(Config.Storage.Dirver, Config.Storage.Config)
	if err != nil {
		log.Fatal(err)
	}
	err = initAdmin(Config.Server.Admin.Account, Config.Server.Admin.Password)
	if err != nil {
		log.Fatal(err)
	}
	// go func() {
	// 	updateContestStanding("./results.xml")
	// 	time.Sleep(time.Duration(1) * time.Second)
	// }()

	gin.SetMode(gin.TestMode)
	GetMainEngine(Config.Server.JWTSecret).Run(":" + Config.Server.Port)
}

func initAdmin(account, password string) error {
	admin := model.User{
		Account:  account,
		Role:     "admin",
		Password: password,
	}
	affected, err := dbEngine.Id(account).AllCols().Update(&admin)
	if err != nil {
		return err
	}
	if affected == 0 {
		_, err = dbEngine.Insert(&admin)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupDB(dirver, config string) error {
	var err error
	dbEngine, err = xorm.NewEngine(dirver, config)
	if err != nil {
		return err
	}
	dbEngine.ShowSQL(true)
	err = dbEngine.Sync2(
		new(model.BallonStatus),
		new(model.User),
		new(model.KVstore),
	)
	return err
}

func getContestStanding(c *gin.Context) {
	var kv model.KVstore
	dbEngine.Where("key = ?", "ContestStanding").Get(&kv)
	c.String(200, string(kv.Value))
}

func getBallonStatus(c *gin.Context) {
	// update ballon status
	var kv model.KVstore
	dbEngine.Where("key = ?", "ContestStanding").Get(&kv)

	var cs model.ContestStanding
	var err = json.Unmarshal(kv.Value, &cs)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internel Server Error"})
	}

	for _, t := range cs.TeamStandings {
		teamKey := t.TeamKey
		for _, p := range t.ProblemSummaryInfos {
			problemIndex := p.Index
			solutionTime := p.SolutionTime
			status := "UnSolved"
			if p.IsFirstSolved {
				status = "FirstBlood"
			} else if p.IsSolved {
				status = "Solved"
			}

			if status == "UnSolved" {
				continue
			}

			var bs model.BallonStatus
			has, _ := dbEngine.Id(core.PK{teamKey, p.Index}).Get(&bs)

			if has {
				dbEngine.Update(&bs, model.BallonStatus{
					Status:       status,
					ProblemIndex: problemIndex,
				})
			} else {
				var user model.User
				dbEngine.Where("TeamKey = ?", teamKey).Get(&user)
				bs = model.BallonStatus{
					TeamKey:      teamKey,
					ProblemIndex: problemIndex,
					TeamName:     user.DisplayName,
					Status:       status,
					SolutionTime: solutionTime,
					IsMarked:     false,
					SeatID:       user.SeatID,
				}
				dbEngine.Insert(&bs)
			}
		}
	}
	var results []model.BallonStatus
	dbEngine.Find(&results)
	data, _ := json.Marshal(results)
	c.String(http.StatusOK, string(data))
}

func patchBallonStatus(c *gin.Context) {
	var req struct {
		TeamKey      string `json:"TeamKey" binding:"required"`
		ProblemIndex int    `json:"ProblemIndex" binding:"required"`
		Action       string `json:"action" binding:"required"`
	}
	if c.BindJSON(&req) == nil {
		log.Println(req)
		if req.Action != "mark" {
			c.JSON(400, gin.H{"message": fmt.Sprintf("No such action=%s", req.Action)})
			return
		}
		var bs model.BallonStatus
		has, err := dbEngine.Id(core.PK{req.TeamKey, req.ProblemIndex}).Get(&bs)
		if err != nil {
			log.Println(err)
		}
		if has {
			affected, err := dbEngine.
				Id(core.PK{req.TeamKey, req.ProblemIndex}).
				Update(&model.BallonStatus{IsMarked: true})
			if err != nil {
				log.Println(err)
			}
			c.JSON(200, gin.H{"affected": affected})
			return
		}
		c.JSON(400, gin.H{"message": fmt.Sprintf("No resource.PK{%s, %d}", req.TeamKey, req.ProblemIndex)})
		return
	}
	c.JSON(400, gin.H{"message": "Invalid Request."})
}

func postPrinter(c *gin.Context) {
	var requestJSON struct {
		PrintContent string `json:"PrintContent" binding:"required"`
	}
	if c.BindJSON(&requestJSON) == nil {
		log.Println("PrintContent:", requestJSON.PrintContent)
		c.JSON(200, gin.H{"status": "OK"})
	} else {
		c.JSON(401, gin.H{"message": "Invalid Request"})
	}
}

func getParticipant(c *gin.Context) {
	var results []model.User
	dbEngine.Find(&results)
	data, err := json.Marshal(results)
	if err != nil {
		c.JSON(501, gin.H{"message": "Internal Server Error."})
	}
	c.String(http.StatusOK, string(data))
}

func postParticipant(c *gin.Context) {
	var err error
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filename := header.Filename
	log.Println("filename : ", filename)

	var allAffected int64
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
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if lineNo == 1 {
			for colNo, title := range A {
				titleMap[title] = colNo
			}
			log.Println(titleMap)
			if _, ok := titleMap["account"]; !ok {
				c.JSON(400, gin.H{"error": fmt.Sprintf("File Format Error. LineNo=%d.", lineNo)})
				return
			}
			if _, ok := titleMap["site"]; !ok {
				c.JSON(400, gin.H{"error": fmt.Sprintf("File Format Error. LineNo=%d.", lineNo)})
				return
			}
		} else {
			if len(A) != len(titleMap) {
				c.JSON(400, gin.H{"error": fmt.Sprintf("File Format Error. LineNo=%d.", lineNo)})
				return
			}
			p := model.User{
				Account: A[titleMap["account"]],
				TeamKey: strings.ToUpper(A[titleMap["site"]] + A[titleMap["account"]]),
			}
			if i, ok := titleMap["password"]; ok {
				p.Password = A[i]
			}
			if i, ok := titleMap["displayname"]; ok {
				p.DisplayName = A[i]
			}
			if i, ok := titleMap["nickname"]; ok {
				p.NickName = A[i]
			}
			if i, ok := titleMap["school"]; ok {
				p.School = A[i]
			}
			if i, ok := titleMap["isstar"]; ok {
				p.IsStar = A[i]
			}
			if i, ok := titleMap["isgirl"]; ok {
				p.IsGirl = A[i]
			}
			if i, ok := titleMap["seatid"]; ok {
				p.SeatID = A[i]
			}
			if i, ok := titleMap["coach"]; ok {
				p.Coach = A[i]
			}
			if i, ok := titleMap["player1"]; ok {
				p.Player1 = A[i]
			}
			if i, ok := titleMap["player2"]; ok {
				p.Player2 = A[i]
			}
			if i, ok := titleMap["player3"]; ok {
				p.Player3 = A[i]
			}
			affected, _ := dbEngine.Id(p.Account).AllCols().Update(&p)
			if affected == 0 {
				affected, _ = dbEngine.Insert(&p)
			}
			allAffected += affected
		}
		lineNo++
	}
	c.JSON(200, gin.H{"affected": allAffected})
}

func genPostAuthenticate(jwtSecret string) func(*gin.Context) {
	return func(c *gin.Context) {
		var requestJSON struct {
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		log.Println(c.Request)
		if c.BindJSON(&requestJSON) == nil {
			log.Println("requestJSON :", requestJSON)
			var user model.User
			has, err := dbEngine.
				Id(requestJSON.Account).
				Cols("Account", "Role", "Password").
				Get(&user)
			if err != nil {
				c.JSON(500, gin.H{"message": "Internal Server Error"})
			}
			log.Println(has, user.Password, requestJSON.Password, user.Password == requestJSON.Password)
			if has && (user.Password == requestJSON.Password) {
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

// JWTAuthMiddleware : JWT Authorization Verification
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt_request.ParseFromRequest(
			c.Request,
			jwt_request.OAuth2Extractor,
			func(token *jwt_lib.Token) (interface{}, error) {
				b := ([]byte(secret))
				return b, nil
			},
		)
		if err != nil {
			c.AbortWithStatus(401)
		}
	}
}

// CORSMiddleware : CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			log.Println("options")
			c.AbortWithStatus(200)
			return
		}
		// c.Next()
	}
}

// GetMainEngine : Main Engine
func GetMainEngine(jwtSecret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(CORSMiddleware())

	api := router.Group("/api")
	{
		postAuthenticate := genPostAuthenticate(jwtSecret)
		api.POST("/authenticate", postAuthenticate)

		authorized := api.Group("/authorized")
		authorized.Use(JWTAuthMiddleware(jwtSecret))
		{
			authorized.GET("/contest-standing", getContestStanding)
			authorized.GET("/ballon-status", getBallonStatus)
			authorized.PATCH("/ballon-status", patchBallonStatus)
			authorized.POST("/printer", postPrinter)
			authorized.GET("/participant", getParticipant)
			authorized.POST("/participant", postParticipant)
		}
	}
	return router
}

func updateContestStanding(resultsXMLPath string) {
	newCS, err := model.ParseResultXML(resultsXMLPath)
	if err != nil {
		log.Println(err)
		return
	}
	b, _ := json.Marshal(newCS)
	newKV := model.KVstore{
		Key:   "ContestStanding",
		Value: b,
	}
	dbEngine.Id("ContestStanding").Update(&newCS)
	var oldKV model.KVstore
	has, err := dbEngine.Id("ContestStanding").Get(&oldKV)
	if err != nil {
		log.Println(err)
		return
	}
	if has {
		// TODO: check diff
		_, err := dbEngine.Id("ContestStanding").Update(&newKV)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err := dbEngine.Insert(newKV)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
