package dbstore

// import (
// 	"encoding/json"

// 	"github.com/shjwudp/ACM-ICPC-api-service/model"
// )

// func (db *dbstore) GetContestStanding() (*model.ContestStanding, error) {
// 	var cs = new(model.ContestStanding)
// 	var content []byte
// 	var err = db.
// 		QueryRow("SELECT content FROM contest_standing ORDER BY timestamp DESC LIMIT 1;").
// 		Scan(&content)
// 	if err != nil {
// 		return cs, err
// 	}
// 	err = json.Unmarshal(content, cs)
// 	return cs, err
// }

// func (db *dbstore) SaveContestStanding(cs *model.ContestStanding) error {
// 	b, err := json.Marshal(cs)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec("UPDATE contest_standing SET content = ?", b)
// }
