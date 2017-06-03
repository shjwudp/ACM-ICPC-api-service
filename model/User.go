package model

// User user struct
type User struct {
	Account     string `xorm:"notnull pk"`
	Password    string
	Role        string
	DisplayName string
	NickName    string
	School      string
	IsStar      string
	IsGirl      string
	SeatID      string
	Coach       string
	Player1     string
	Player2     string
	Player3     string
	Site        string
	TeamKey     string `xorm:"unique"`
}
