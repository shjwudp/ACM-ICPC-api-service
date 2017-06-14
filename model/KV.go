package model

import (
	"bytes"
	"crypto/md5"
	"database/sql"
)

// KV simple kv
type KV struct {
	Key      string `db:"key"`
	Value    []byte `db:"value"`
	CheckSum []byte `db:"check_sum"`
}

// GetKV get KV by Key
func (db *DB) GetKV(key string) (*KV, error) {
	kv := new(KV)
	err := db.Get(kv, "SELECT * FROM kv WHERE key=$1", key)
	return kv, err
}

// GetKVCheckSum get checksum of KV.Value
func (db *DB) GetKVCheckSum(key string) (*KV, error) {
	kv := &KV{Key: key}
	selectSQL := "SELECT check_sum FROM kv WHERE key=$1"
	err := db.QueryRow(selectSQL, key).
		Scan(&kv.CheckSum)
	return kv, err
}

// SaveKV save kv in db
func (db *DB) SaveKV(kv KV) error {
	md5Array := md5.Sum(kv.Value)
	kv.CheckSum = md5Array[0:len(md5Array)]
	last, err := db.GetKVCheckSum(kv.Key)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows || !bytes.Equal(last.CheckSum, kv.CheckSum) {
		insertSQL := "INSERT OR REPLACE INTO kv VALUES($1, $2, $3)"
		_, err = db.Exec(insertSQL, kv.Key, kv.Value, kv.CheckSum)
		return err
	}
	return nil
}
