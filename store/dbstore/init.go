package dbstore

// import (
// 	"database/sql"

// 	_ "github.com/mattn/go-sqlite3"
// )

// // CreateAllTable create all tables;
// func CreateAllTable(db *sql.DB) error {
// 	var err error
// 	dbEngine, err = xorm.NewEngine("sqlite3", "./sqlite3.db")
// 	if err != nil {
// 		return err
// 	}
// 	err = dbEngine.Sync2(
// 		new(models.BallonStatus),
// 		new(models.Participant),
// 		new(models.Problem),
// 		new(models.StandingsHeader),
// 		new(models.ProblemSummaryInfo),
// 		new(models.TeamStanding),
// 		new(models.ContestStanding),
// 	)
// }

// var createContestStandingTable = `
// CREATE TABLE IF NOT EXISTS contest_standing (
// 	content BLOB,
// 	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
// )
// `
// var createUsersTable = `
// CREATE TABLE IF NOT EXISTS users (
//     account TEXT PRIMARY KEY NOT NULL,
//     password TEXT NOT NULL,
//     role TEXT NOT NULL DEFAULT 'normal',
//     display_name TEXT NULL,
//     nick_name TEXT NULL,
//     school TEXT NULL,
//     is_star TEXT NULL,
//     is_girl TEXT NULL,
//     seat_id TEXT NULL,
//     coach TEXT NULL,
//     player1 TEXT NULL,
//     player2 TEXT NULL,
//     player3 TEXT NULL
// )
// `
