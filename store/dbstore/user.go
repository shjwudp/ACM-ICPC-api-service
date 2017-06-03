package dbstore

// import (
// 	"encoding/json"

// 	"github.com/shjwudp/ACM-ICPC-api-service/model"
// )

// func (db *DB) GetUser(user *model.User) error {
// 	var sql = `
// 	SELECT password, role, display_name, nick_name,
// 		school, is_star, is_girl, seat_id, coach,
// 		player1, player2, player3, site, team_key
// 	FROM user
// 	WHERE account = ?
// 	`
// 	var err = db.QueryRow(sql,
// 		 &user.Password,
// 		 &user.Role,
// 		 &user.DisplayName,

// 		 )
// 	var user User{account:account}
// 	var selectSql = `
// 	SELECT (password, role, display_name, nick_name, school, is_star, is_girl, seat_id, coach, player1, player2, player3)
// 	FROM user
// 	WHERE account = ?
// 	`
// 	var err = db.
// 		QueryRow(selectSql, account).
// 		Scan(&user)
// 	return user, err
// }

// func (db *dbstore) SaveUser(user *model.User) error {
// 	var user User{}
// }

// func (db *dbstore) GetUserList() ([]*model.User, error) {
// 	b, err := json.Marshal(cs)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec("UPDATE contest_standing SET content = ?", b)
// }
