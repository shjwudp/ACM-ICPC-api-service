package store

import (
	"github.com/shjwudp/ACM-ICPC-api-service/model"
)

// Store interface
type Store interface {
	// SaveContestStanding
	SaveContestStanding(*model.ContestStanding) error

	// GetContestStanding get lastest
	GetContestStanding() (*model.ContestStanding, error)

	// GetUser get user by pk account.
	GetUser(string) (*model.User, error)

	SaveUser(*model.User)

	// GetUserList gets a list of all users.
	GetUserList() ([]*model.User, error)

	// GetBallonStatusList gets a list of all ballonstatus.
	GetBallonStatusList() ([]*model.BallonStatus, error)
}
