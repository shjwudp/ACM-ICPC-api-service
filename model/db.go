package model

import (
	// "database/sql"
	"github.com/jmoiron/sqlx"
)

// Datastore interface for model operate
type Datastore interface {
	// GetUserAccount get user by Account
	GetUserAccount(string) (*User, error)

	// GetUserTeamKey get user by TeamKey
	GetUserTeamKey(string) (*User, error)

	SaveUser(User) error

	// AllUser get all of user
	AllUser() ([]User, error)

	// GetKV get KV by Key
	GetKV(string) (*KV, error)

	SaveKV(KV) error

	// GetBallonStatus get BallonStatus by TeamKey AND ProblemIndex
	GetBallonStatus(teamKey string, pid int) (*BallonStatus, error)

	SaveBallonStatus(BallonStatus) error

	SavePrintCode(account, code string) (*PrintCode, error)

	// GetPrintCode get PrintCode by ID
	GetPrintCode(id int64) (*PrintCode, error)

	// UpdatePrintCode update PrintCode(Account, Code, IsDown) WHERE ID = p.ID
	UpdatePrintCode(p PrintCode) error

	// GetContestEventSeq get ContestEvent by TeamKey
	GetContestEventSeq(teamKey string) ([]ContestEvent, error)

	// GetLatestContestEvent get latest ContestEvent by TeamKey & ProblemIndex
	GetLatestContestEvent(teamKey string, problemIndex int) (*ContestEvent, error)

	// SaveContestEvent save ContestEvent in db
	SaveContestEvent(ce ContestEvent) error
}

// DB is an implementation of a store.Store built on top
type DB struct {
	*sqlx.DB
}

// OpenDB creates a database connection for the given driver and datasource
// and returns a new Store.
func OpenDB(conf StorageConfiguration) (*DB, error) {
	db, err := sqlx.Connect(conf.Dirver, conf.Addr)
	// db, err := sql.Open(driver, config)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.MustExec(migrateSQL)
	return &DB{db}, nil
}

// OpenDBTest open a temporary DB for test
func OpenDBTest() (*DB, error) {
	return OpenDB(StorageConfiguration{
		Dirver:       "sqlite3",
		Addr:         ":memory:",
		MaxIdleConns: 1000,
		MaxOpenConns: 600,
	})
}

var migrateSQL = `
PRAGMA read_uncommitted = 1;
CREATE TABLE IF NOT EXISTS ballon_status (
    team_key TEXT NOT NULL,
    problem_index INTEGER NOT NULL,
    is_marked INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (team_key, problem_index)
);
CREATE TABLE IF NOT EXISTS kv (
    key TEXT PRIMARY KEY NOT NULL,
    value BLOB NOT NULL DEFAULT '',
	check_sum TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS user (
    account TEXT PRIMARY KEY NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'normal',
    display_name TEXT NOT NULL DEFAULT '',
    nick_name TEXT NOT NULL DEFAULT '',
    school TEXT NOT NULL DEFAULT '',
    is_star INTEGER NOT NULL DEFAULT 0,
    is_girl INTEGER NOT NULL DEFAULT 0,
    seat_id TEXT NOT NULL DEFAULT '',
    coach TEXT NOT NULL DEFAULT '',
    player1 TEXT NOT NULL DEFAULT '',
    player2 TEXT NOT NULL DEFAULT '',
    player3 TEXT NOT NULL DEFAULT '',
    site TEXT NOT NULL DEFAULT '',
    team_key TEXT NOT NULL UNIQUE DEFAULT ''
);
CREATE TABLE IF NOT EXISTS print_code (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	account TEXT NOT NULL,
	code TEXT NOT NULL,
	is_done INTEGER NOT NULL DEFAULT 0,
	create_time INTEGER not null default (strftime('%s','now'))
);
CREATE TABLE IF NOT EXISTS contest_event (
	time_stamp INTEGER NOT NULL,
	team_key TEXT NOT NULL,
	problem_index INTEGER NOT NULL,
	attempts INTEGER NOT NULL,
	is_solved INTEGER NOT NULL,
	points INTEGER NOT NULL,
	PRIMARY KEY (time_stamp, team_key, problem_index)
);
CREATE INDEX IF NOT EXISTS contest_event__time_stamp__index ON contest_event(time_stamp);
CREATE INDEX IF NOT EXISTS contest_event__team_key__index ON contest_event(team_key);
`
