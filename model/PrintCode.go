package model

// PrintCode ACM-ICPC PrintCode
type PrintCode struct {
	ID         int64  `db:"id"`
	Account    string `db:"account"`
	Code       string `db:"code"`
	IsDone     bool   `db:"is_done"`
	CreateTime int64  `db:"create_time"`
}

func (db *DB) SavePrintCode(account, code string) (*PrintCode, error) {
	var saveSQL = `
	INSERT INTO print_code ( account, code )
	VALUES ( $1, $2 )
	`
	res, err := db.Exec(saveSQL, account, code)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return db.GetPrintCode(id)
}

// GetPrintCode get PrintCode by ID
func (db *DB) GetPrintCode(id int64) (*PrintCode, error) {
	p := new(PrintCode)
	err := db.Get(p, "SELECT * FROM print_code WHERE id = $1", id)
	return p, err
}

// UpdatePrintCode update PrintCode(Account, Code, IsDown) WHERE ID = p.ID
func (db *DB) UpdatePrintCode(p PrintCode) error {
	var updateSQL = `
	UPDATE print_code
	SET account = :account, code = :code, is_done = :is_done
	WHERE id = :id
	`
	_, err := db.NamedExec(updateSQL, &p)
	return err
}
