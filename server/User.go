package server

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"io"
	"log"
	"net/http"
	"strings"
)

// AllUser list all of User
func (env *Env) AllUser(c *gin.Context) {
	users, err := env.db.AllUser()
	if err != nil {
		errMsg := fmt.Sprint("List User failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	c.JSON(http.StatusOK, users)
}

// PostUserList post userlist and save them.
func (env *Env) PostUserList(c *gin.Context) {
	var err error
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		errMsg := fmt.Sprint("Get uploadFile failed with", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
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
			c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
			return
		}
		if lineNo == 1 {
			for colNo, title := range A {
				titleMap[title] = colNo
			}
			log.Println(titleMap)
			if _, ok := titleMap["account"]; !ok {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
				return
			}
			if _, ok := titleMap["site"]; !ok {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
				return
			}
		} else {
			if len(A) != len(titleMap) {
				errMsg := fmt.Sprintf("Read uploadFile:%d failed with %s", lineNo, err)
				c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
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
				if A[i] == "true" {
					u.IsStar = true
				} else {
					u.IsStar = false
				}
			}
			if i, ok := titleMap["isgirl"]; ok {
				if A[i] == "true" {
					u.IsGirl = true
				} else {
					u.IsGirl = false
				}
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
				c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
				return
			}
		}
		lineNo++
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
