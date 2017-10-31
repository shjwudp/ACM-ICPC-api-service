package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"log"
	"net/http"
)

// GetContestInfo get ContestInfo from db
func (env *Env) GetContestInfo(c *gin.Context) {
	kv, err := env.db.GetKV("ContestInfo")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestInfo failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	var ci model.ContestInfo
	err = json.Unmarshal(kv.Value, &ci)
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestInfo failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	c.JSON(http.StatusOK, ci)
}

// SaveContestInfo save ContestInfo in db
func (env *Env) SaveContestInfo(c *gin.Context) {
	var requestJSON model.ContestInfo
	log.Println(c.Request)
	err := c.BindJSON(&requestJSON)
	if err != nil {
		log.Println(err)
		errMsg := fmt.Errorf("BindJSON failed with %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}
	b, err := json.Marshal(requestJSON)
	if err != nil {
		errMsg := fmt.Errorf("json.Marshal failed with %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	err = env.db.SaveKV(model.KV{Key: "ContestInfo", Value: b})
	if err != nil {
		errMsg := fmt.Errorf("db.SaveKV failed with %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
