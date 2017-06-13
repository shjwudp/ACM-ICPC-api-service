package model

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_ContestEvent(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Fatal(err)
	}
	eventList := []ContestEvent{
		ContestEvent{
			TimeStamp:    1234,
			TeamKey:      "team1",
			ProblemIndex: 1,
			Attempts:     1,
			IsSolved:     false,
			Points:       0,
		},
		ContestEvent{
			TimeStamp:    1235,
			TeamKey:      "team1",
			ProblemIndex: 1,
			Attempts:     1,
			IsSolved:     false,
			Points:       0,
		},
		ContestEvent{
			TimeStamp:    1236,
			TeamKey:      "team1",
			ProblemIndex: 1,
			Attempts:     2,
			IsSolved:     true,
			Points:       21,
		},
	}
	for _, ce := range eventList {
		err = db.SaveContestEvent(ce)
		if err != nil {
			t.Fatal(err)
		}
	}

	ceList, err := db.GetContestEventSeq("team1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ceList) != 2 {
		t.Fatal("len(ceList) must be 2")
	}
}
