package model

import (
	"bytes"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_KV(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Fatal(err)
	}
	kvList := []KV{
		KV{
			Key:   "abc",
			Value: []byte("value A"),
		},
		KV{
			Key:   "abc",
			Value: []byte("value B"),
		},
	}

	err = db.SaveKV(kvList[0])
	if err != nil {
		t.Fatal(err)
	}
	err = db.SaveKV(kvList[0])
	if err != nil {
		t.Fatal(err)
	}
	err = db.SaveKV(kvList[1])
	if err != nil {
		t.Fatal(err)
	}

	q, err := db.GetKV(kvList[1].Key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(q.Value, kvList[1].Value) {
		t.Fatal("Assert q.Value == kv2.Value")
	}
}
