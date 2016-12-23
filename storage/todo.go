package storage

type Todo struct {
	Id       int64 `db:"id" json:"id"`
	Item     int64 `db:"item" json:"item"`
	TicketId int64 `db:"ticket_id" json:"ticket_id"`
	Active   bool  `db:"active" json:"active"`
}

func (adb *AppDB) CreateTodo(t *Todo) (int64, error) {
	result, err := adb.db.NamedExec("INSERT INTO todo (id, item, ticket_id, active) VALUES (:id, :item, :ticket_id, :active)", t)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
