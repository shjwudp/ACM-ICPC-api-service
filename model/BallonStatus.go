package model

type BallonStatus struct {
	TeamKey      string `xorm:"notnull pk"`
	ProblemIndex int    `xorm:"notnull pk"`
	SolutionTime int
	TeamName     string
	SeatID       string
	Status       string
	IsMarked     bool
}
