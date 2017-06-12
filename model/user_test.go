package model

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Test_User(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Error(err)
	}
	user := User{
		Account:     "hasaki",
		Password:    "haha",
		Role:        "good man",
		DisplayName: "nico",
		NickName:    "niconiconi",
		School:      "fofofo",
		IsStar:      "0",
		IsGirl:      "1",
		Coach:       "MC",
		Player1:     "Player1",
		Player2:     "Player2",
		Player3:     "Player3",
		Site:        "5",
		TeamKey:     "5HASAKI",
	}
	user.TeamKey = strings.ToUpper(user.Site + user.Account)

	err = db.SaveUser(user)
	if err != nil {
		t.Error(err)
	}
	user.IsStar = "1"
	err = db.SaveUser(user)
	if err != nil {
		t.Error(err)
	}

	getUser, err := db.GetUserAccount("hasak")
	if err != sql.ErrNoRows {
		t.Error(err)
	}
	getUser, err = db.GetUserAccount("hasaki")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(&user, getUser) {
		t.Error("Assert q.Value == kv2.Value")
	}
	getUser, err = db.GetUserTeamKey("5HASAKI")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(&user, getUser) {
		t.Error("Assert q.Value == kv2.Value")
	}

	userList, err := db.ListUser()
	if err != nil {
		t.Error(err)
	}
	log.Println(len(userList), userList)
	if len(userList) != 1 {
		t.Error("Assert len(userList) == 1")
	}

	data, err := json.Marshal(userList)
	log.Println(string(data))
}
