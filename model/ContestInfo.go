package model

import (
	"time"
)

// ContestInfo describe contest
type ContestInfo struct {
	StartTime      time.Time     `json:"StartTime" binding:"required"`
	GoldMedalNum   int           `json:"GoldMedalNum" binding:"required"`
	SilverMedalNum int           `json:"SilverMedalNum" binding:"required"`
	BronzeMedalNum int           `json:"BronzeMedalNum" binding:"required"`
	Duration       time.Duration `json:"Duration" binding:"required"`
}
