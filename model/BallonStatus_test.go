package model

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_BallonStatus(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Fatal(err)
	}
	b := BallonStatus{
		TeamKey:      "1TEAM1",
		ProblemIndex: 3,
		IsMarked:     true,
	}
	err = db.SaveBallonStatus(b)
	if err != nil {
		t.Fatal(err)
	}
	b.IsMarked = true
	err = db.SaveBallonStatus(b)
	if err != nil {
		t.Fatal(err)
	}
	saveB, err := db.GetBallonStatus(b.TeamKey, b.ProblemIndex)
	if saveB.IsMarked != b.IsMarked {
		t.Fatal("Assert saveB.IsMarked == b.IsMarked")
	}
}
