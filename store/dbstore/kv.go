package dbstore

import (
	"errors"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
)

func (db *DB) GetKV(kv *model.KV) error {
	if kv == nil {
		return errors.New("Unexcepted nil")
	}
	err := db.Get(kv, "SELECT * FROM kv WHERE key=$1", kv.Key)
	// err := db.QueryRow("SELECT value FROM kv WHERE key=$1", kv.Key).Scan(&kv.Value)
	// err := db.Get(&kv, "SELECT * FROM kv WHERE key=$1", kv.Key)
	return err
}

func (db *DB) SaveKV(kv *model.KV) error {
	if kv == nil {
		return errors.New("Unexcepted nil")
	}
	_, err := db.Exec("INSERT OR REPLACE INTO kv VALUES($1, $2)", kv.Key, kv.Value)
	return err
}
