package storage

type Todo struct {
	Id       int64  `db:"id" json:"id"`
	Idx      int64  `db:"idx" json:"idx"`
	Item     string `db:"item" json:"item"`
	TicketId int64  `db:"ticket_id" json:"ticket_id"`
	Active   bool   `db:"active" json:"active"`
}

func (adb *AppDB) CreateTodo(t *Todo) (int64, error) {
	result, err := adb.db.NamedExec("INSERT INTO todo (idx, item, ticket_id, active) VALUES (:idx, :item, :ticket_id, :active)", t)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
