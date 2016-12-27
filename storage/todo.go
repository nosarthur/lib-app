package storage

type Todo struct {
	Id       int64  `db:"id" json:"id"`
	Idx      int64  `db:"idx" json:"idx"`
	Item     string `db:"item" json:"item"`
	TicketId string `db:"ticket_id" json:"ticket_id"`
	Done     bool   `db:"done" json:"done"`
}

func (adb *AppDB) getTodoCount(ticket_id string) (int64, error) {
	var count int64
	err := adb.db.QueryRow("SELECT COUNT(*) FROM todo WHERE ticket_id=$1", ticket_id).Scan(&count)
	return count, err
}

func (adb *AppDB) CreateTodo(t Todo) error {
	if _, err := adb.ReadTicket(t.TicketId); err != nil {
		return err
	}
	count, err := adb.getTodoCount(t.TicketId)
	if err != nil {
		return err
	}
	t.Idx = count + 1
	_, err = adb.db.NamedExec("INSERT INTO todo (idx, item, ticket_id, done) VALUES (:idx, :item, :ticket_id, :done)", &t)
	return err
}

func (adb *AppDB) ReadTodo(ticket_id string, idx int64) (Todo, error) {
	t := Todo{}
	query := `SELECT * from todo WHERE ticket_id=$1 AND idx=$2`
	err := adb.db.Get(&t, query, ticket_id, idx)
	return t, err
}

func (adb *AppDB) ReadTodos(ticket_id string) ([]Todo, error) {
	query := `SELECT * from todo WHERE ticket_id=$1`
	rows, err := adb.db.Queryx(query, ticket_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		value := Todo{}
		if err = rows.StructScan(&value); err != nil {
			return nil, err
		}
		todos = append(todos, value)
	}
	return todos, rows.Err()
}

func (adb *AppDB) UpdateTodo(t Todo) error {
	_, err := adb.db.NamedExec(`UPDATE todo SET item=:item, done=:done WHERE id=:id;`, &t)
	return err
}

func (adb *AppDB) DeleteTodo(t Todo) error {
	_, err := adb.db.NamedExec(`DELETE FROM todo WHERE id=:id AND idx=:idx;`, &t)
	return err
}
