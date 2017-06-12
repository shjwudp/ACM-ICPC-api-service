package model

// KV simple kv
type KV struct {
	Key   string `xorm:"notnull pk"`
	Value []byte
}

// GetKV get KV by Key
func (db *DB) GetKV(key string) (*KV, error) {
	kv := new(KV)
	err := db.Get(kv, "SELECT * FROM kv WHERE key=$1", key)
	// err := db.QueryRow("SELECT value FROM kv WHERE key=$1", kv.Key).Scan(&kv.Value)
	// err := db.Get(&kv, "SELECT * FROM kv WHERE key=$1", kv.Key)
	return kv, err
}

// SaveKV save kv in db
func (db *DB) SaveKV(kv KV) error {
	_, err := db.Exec("INSERT OR REPLACE INTO kv VALUES($1, $2)", kv.Key, kv.Value)
	return err
}
