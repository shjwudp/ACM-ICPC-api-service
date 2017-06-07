package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetContestStanding get kv.Value ContestStanding
func (env *Env) GetContestStanding(c *gin.Context) {
	kv, err := env.db.GetKV("ContestStanding")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
	}
	c.String(200, string(kv.Value))
}
