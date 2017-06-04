package model

// BallonStatus ACM-ICPC ballon status
// Status must in ("UnSolved" "FirstBlood" "Solved")
type BallonStatus struct {
	TeamKey      string `db:"team_key" xorm:"notnull pk"`
	ProblemIndex int    `db:"problem_index" xorm:"notnull pk"`
	IsMarked     bool   `db:"is_marked"`
}

// GetBallonStatus get BallonStatus by TeamKey AND ProblemIndex
func (db *DB) GetBallonStatus(teamKey string, pid int) (*BallonStatus, error) {
	b := new(BallonStatus)
	var selectSQL = "SELECT * FROM ballon_status WHERE team_key = $1 AND problem_index = $2"
	err := db.Get(b, selectSQL, teamKey, pid)
	return b, err
}

func (db *DB) SaveBallonStatus(b BallonStatus) error {
	var saveSQL = `
	INSERT OR REPLACE INTO ballon_status
	( team_key, problem_index, is_marked )
	VALUES
	( :team_key, :problem_index, :is_marked )
	`
	_, err := db.NamedExec(saveSQL, &b)
	return err
}
