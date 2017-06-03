package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/shjwudp/ACM-ICPC-api-service/model"

	jwt_lib "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gin-gonic/gin.v1"
	"strings"
)

var dbEngine *xorm.Engine

func initAdmin() error {
	admin := model.Participant{
		Account:  "admin",
		Role:     "admin",
		Password: "ElPsyCongroo",
	}
	affected, err := dbEngine.Id("admin").AllCols().Update(&admin)
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
	err = dbEngine.Sync2(
		new(model.BallonStatus),
		new(model.Participant),
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
			has, _ := dbEngine.Id(core.PK{teamKey, problemIndex}).Get(&bs)

			if has {
				bs.Status = status
				bs.SolutionTime = solutionTime
				dbEngine.
					Id(core.PK{teamKey, problemIndex}).
					Cols("Status", "SolutionTime").
					Update(bs)
			} else {
				var user model.user
				dbEngine.Where("TeamKey = ?", teamKey).Get(&user)
				bs = model.BallonStatus{
					TeamKey: teamKey,
					ProblemIndex: problemIndex,
					TeamName:     teamName,
					Status:       status,
					SolutionTime: solutionTime,
					Read:     false,
					SeatId: user.SeatId
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
	teamKey := c.Query("TeamKey")
	problemIndex := c.Query("ProblemIndex")
	action := c.Query("action")
	if teamKey == "" || problemIndex == "" || action != "read" {
		c.JSON(400, gin.H{"message": "Method Not Allowed."})
		return
	}

	var bs model.BallonStatus
	has, _ := dbEngine.Id(core.PK{teamKey, problemIndex}).Get(&bs)
	if has {
		bs.Read = true
		affected, _ := dbEngine.Cols("Read").Update(&bs)
		c.JSON(200, gin.H{"affected": affected})
	} else {
		c.JSON(400, gin.H{"message": fmt.Sprintf("No resource.PK{%s, %s}", teamKey, problemIndex)})
	}
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
	var results []model.Participant
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
			p := model.Participant{
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
				p.SeatId = A[i]
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

func genPostAuthenticate(jwtSecret string) {
	return func(c *gin.Context) {
		var requestJSON struct {
			Account  string `json:"account" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		log.Println(c.Request)
		if c.BindJSON(&requestJSON) == nil {
			var user model.Participant
			has, err := dbEngine.
				Id(requestJSON.Account).
				Cols("Account", "Role", "Password").
				Get(&user)
			if err != nil {
				c.JSON(500, gin.H{"message": "Internal Server Error"})
			}
			if has && user.Password == requestJSON.Password {
				log.Println(user)
				// Create the token
				token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
				// Set some claims
				token.Claims = jwt_lib.MapClaims{
					"Account": user.Account,
					"Role":    user.Role,
					"exp":     time.Now().Add(time.Hour * 12).Unix(),
				}
				log.Println("jwtSecret : ", jwtSecret)
				// Sign and get the complete encoded token as a string
				tokenString, err := token.SignedString([]byte(jwtSecret))
				if err != nil {
					c.JSON(500, gin.H{"message": "Could not generate token"})
				}
				c.JSON(200, gin.H{"token": tokenString})
				return
			}
			c.JSON(401, gin.H{"message": "Wrong account or password"})
		} else {
			c.JSON(401, gin.H{"message": "Invalid Request"})
		}
		log.Printf("%v\n", requestJSON)
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
			// c.AbortWithStatus(401)
		}
	}
}

// CORSMiddleware : CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
	var oldCS model.ContestStanding
	dbEngine.Get(&oldCS)
	has, err := dbEngine.Get(&oldCS)
	if err != nil {
		log.Println(err)
		return
	}
	if has {
		_, err = dbEngine.Update(&oldCS, &newCS)
	} else {
		_, err = dbEngine.Insert(&newCS)
	}
	log.Println(err)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	err = initAdmin()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		updateContestStanding("./results.xml")
		time.Sleep(time.Duration(1) * time.Second)
	}()

	gin.SetMode(gin.TestMode)
	GetMainEngine().Run(":8080")
}
