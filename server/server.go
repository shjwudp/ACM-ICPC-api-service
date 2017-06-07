package server

import (
	"github.com/shjwudp/ACM-ICPC-api-service/model"
)

// Env is all handlers common environment
type Env struct {
	db         model.Datastore
	printQueue chan int64
	jwtSecret  string
}

// NewEnv return a env
func NewEnv(db model.Datastore, printQueueLength int, JWTSecret string) *Env {
	return &Env{
		db:         db,
		printQueue: make(chan int64, printQueueLength),
		jwtSecret:  JWTSecret,
	}
}
