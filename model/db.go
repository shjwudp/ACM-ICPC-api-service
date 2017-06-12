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

	ListUser() ([]User, error)

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
}

// DB is an implementation of a store.Store built on top
type DB struct {
	*sqlx.DB
}

// OpenDB creates a database connection for the given driver and datasource
// and returns a new Store.
func OpenDB(driver, config string) (*DB, error) {
	db, err := sqlx.Connect(driver, config)
	// db, err := sql.Open(driver, config)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(400)
	db.SetMaxOpenConns(400)
	db.MustExec(migrateSQL)
	return &DB{db}, nil
}

// OpenDBTest open a temporary DB for test
func OpenDBTest() (*DB, error) {
	return OpenDB("sqlite3", ":memory:")
}

var migrateSQL = `
PRAGMA cache_size = 400;
CREATE TABLE IF NOT EXISTS ballon_status (
    team_key TEXT NOT NULL,
    problem_index INTEGER NOT NULL,
    is_marked INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (team_key, problem_index)
);
CREATE TABLE IF NOT EXISTS kv (
    key TEXT PRIMARY KEY NOT NULL,
    value BLOB NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS user (
    account TEXT PRIMARY KEY NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'normal',
    display_name TEXT NOT NULL DEFAULT '',
    nick_name TEXT NOT NULL DEFAULT '',
    school TEXT NOT NULL DEFAULT '',
    is_star TEXT NOT NULL DEFAULT '',
    is_girl TEXT NOT NULL DEFAULT '',
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
`
