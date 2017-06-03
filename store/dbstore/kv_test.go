package dbstore

import (
	"bytes"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_kv(t *testing.T) {
	db, err := NewDB("sqlite3", "./test.db")
	if err != nil {
		t.Error(err)
	}
	kv1 := model.KV{
		Key:   "abc",
		Value: []byte("value A"),
	}
	kv2 := model.KV{
		Key:   "abc",
		Value: []byte("value B"),
	}

	err = db.SaveKV(&kv1)
	if err != nil {
		t.Error(err)
	}
	err = db.SaveKV(&kv2)
	if err != nil {
		t.Error(err)
	}

	var q = model.KV{Key: "abc"}
	err = db.GetKV(&q)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(q.Value, kv2.Value) {
		t.Error("Assert q.Value == kv2.Value")
	}
}
