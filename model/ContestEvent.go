package model

import (
	"database/sql"
	"reflect"
)

// ContestEvent record basic contest event
type ContestEvent struct {
	TimeStamp    int64  `db:"time_stamp"`
	TeamKey      string `db:"team_key"`
	ProblemIndex int    `db:"problem_index"`
	Attempts     int    `db:"attempts"`
	IsSolved     bool   `db:"is_solved"`
	Points       int    `db:"points"`
}

// GetContestEventSeq get ContestEvent by TeamKey
func (db *DB) GetContestEventSeq(teamKey string) ([]ContestEvent, error) {
	ce := []ContestEvent{}
	selectSQL := `
	SELECT * FROM contest_event 
	WHERE team_key=$1
	ORDER BY time_stamp ASC
	`
	err := db.Select(&ce, selectSQL, teamKey)
	return ce, err
}

// GetLatestContestEvent get latest ContestEvent by TeamKey & ProblemIndex
func (db *DB) GetLatestContestEvent(teamKey string, problemIndex int) (*ContestEvent, error) {
	ce := new(ContestEvent)
	selectSQL := `
	SELECT * FROM contest_event 
	WHERE team_key=$1 AND problem_index=$2 
	ORDER BY time_stamp DESC LIMIT 1
	`
	err := db.Get(ce, selectSQL, teamKey, problemIndex)
	return ce, err
}

// SaveContestEvent save ContestEvent in db
func (db *DB) SaveContestEvent(ce ContestEvent) error {
	last, err := db.GetLatestContestEvent(ce.TeamKey, ce.ProblemIndex)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	var needAppend bool
	if err == sql.ErrNoRows {
		needAppend = true
	} else {
		last.TimeStamp = ce.TimeStamp
		needAppend = !reflect.DeepEqual(*last, ce)
	}
	if !needAppend {
		return nil
	}

	var saveSQL = `
	INSERT OR REPLACE INTO contest_event
	( time_stamp, team_key, problem_index,
	attempts, is_solved, points )
	VALUES
	( :time_stamp, :team_key, :problem_index,
	:attempts, :is_solved, :points )
	`
	_, err = db.NamedExec(saveSQL, &ce)
	return err
}
