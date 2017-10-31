package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"log"
	"net/http"
)

// ListBallonStatus List all of BallonStatus
func (env *Env) ListBallonStatus(c *gin.Context) {
	// update BallonStatus first
	kv, err := env.db.GetKV("ContestStanding")
	if err != nil {
		var errMsg = fmt.Sprint("Get ContestStanding failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}

	var cs model.ContestStanding
	err = json.Unmarshal(kv.Value, &cs)
	if err != nil {
		var errMsg = fmt.Sprint("Json Unmarshal failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
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
				c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
				return
			}
			results = append(results, map[string]interface{}{
				"TeamName":      t.TeamName,
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
	c.JSON(http.StatusOK, results)
}

// PatchBallonStatus modify one BallonStatus
func (env *Env) PatchBallonStatus(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
			return
		}
		err := env.db.SaveBallonStatus(model.BallonStatus{
			TeamKey:      req.TeamKey,
			ProblemIndex: req.ProblemIndex,
			IsMarked:     true,
		})
		if err != nil {
			errMsg := fmt.Sprintf("No such action=%s", req.Action)
			c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprint("BindJSON failed with", err)})
}
