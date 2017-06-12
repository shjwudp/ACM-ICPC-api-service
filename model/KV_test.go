package model

import (
	"bytes"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_KV(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Error(err)
	}
	kv1 := KV{
		Key:   "abc",
		Value: []byte("value A"),
	}
	kv2 := KV{
		Key:   "abc",
		Value: []byte("value B"),
	}

	err = db.SaveKV(kv1)
	if err != nil {
		t.Error(err)
	}
	err = db.SaveKV(kv2)
	if err != nil {
		t.Error(err)
	}

	q, err := db.GetKV("abc")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(q.Value, kv2.Value) {
		t.Error("Assert q.Value == kv2.Value")
	}
}
