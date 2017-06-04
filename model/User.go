package model

// User user struct
type User struct {
	Account     string `db:"account" xorm:"notnull pk"`
	Password    string `db:"password"`
	Role        string `db:"role"`
	DisplayName string `db:"display_name"`
	NickName    string `db:"nick_name"`
	School      string `db:"school"`
	IsStar      string `db:"is_star"`
	IsGirl      string `db:"is_girl"`
	SeatID      string `db:"seat_id"`
	Coach       string `db:"coach"`
	Player1     string `db:"player1"`
	Player2     string `db:"player2"`
	Player3     string `db:"player3"`
	Site        string `db:"site"`
	TeamKey     string `db:"team_key" xorm:"unique"`
}

// GetUserAccount get user by Account
func (db *DB) GetUserAccount(account string) (*User, error) {
	user := new(User)
	err := db.Get(user, "SELECT * FROM user WHERE account = $1", account)
	return user, err
}

// GetUserTeamKey get user by TeamKey
func (db *DB) GetUserTeamKey(teamKey string) (*User, error) {
	user := new(User)
	err := db.Get(user, "SELECT * FROM user WHERE team_key = $1", teamKey)
	return user, err
}

func (db *DB) ListUser() ([]User, error) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM user")
	return users, err
}

func (db *DB) SaveUser(user User) error {
	var saveSQL = `
	INSERT OR REPLACE INTO user
	( account, password, role, display_name, nick_name,
	school, is_star, is_girl, seat_id, coach, player1, 
	player2, player3, site, team_key )
	VALUES
	( :account, :password, :role, :display_name, :nick_name,
	:school, :is_star, :is_girl, :seat_id, :coach, :player1, 
	:player2, :player3, :site, :team_key )
	`
	_, err := db.NamedExec(saveSQL, &user)
	return err
}
