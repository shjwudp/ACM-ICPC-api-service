package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"net/http"
)

// GetContestStanding get kv.Value ContestStanding
func (env *Env) GetContestStanding(c *gin.Context) {
	kv, err := env.db.GetKV("ContestStanding")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	var cs model.ContestStanding
	err = json.Unmarshal(kv.Value, &cs)
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	c.JSON(http.StatusOK, cs)
}
