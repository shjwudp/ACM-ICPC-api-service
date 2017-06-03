package model

// KV simple kv
type KV struct {
	Key   string `xorm:"notnull pk"`
	Value []byte
}
