package dbstore

import (
	// "database/sql"
	"github.com/jmoiron/sqlx"
)

// DB is an implementation of a store.Store built on top
type DB struct {
	*sqlx.DB
}

// NewDB creates a database connection for the given driver and datasource
// and returns a new Store.
func NewDB(driver, config string) (*DB, error) {
	db, err := sqlx.Connect(driver, config)
	// db, err := sql.Open(driver, config)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
