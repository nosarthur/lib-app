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
	err := adb.db.QueryRow("SELECT COUNT(*) FROM todo WHERE ticket_id=?", ticket_id).Scan(&count)
	return count, err
}

func (adb *AppDB) CreateTodo(t Todo) (int64, error) {
	_, err := adb.ReadTicket(t.TicketId)
	if err != nil {
		return 0, err
	}
	count, err := adb.getTodoCount(t.TicketId)
	if err != nil {
		return 0, err
	}
	t.Idx = count + 1
	result, err := adb.db.NamedExec("INSERT INTO todo (idx, item, ticket_id, done) VALUES (:idx, :item, :ticket_id, :done)", &t)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (adb *AppDB) ReadTodo(ticket_id string, idx int64) (Todo, error) {
	t := Todo{}
	query := `SELECT * from todo WHERE ticket_id=? AND idx=?`
	err := adb.db.Get(&t, query, ticket_id, idx)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (adb *AppDB) ReadTodos(ticket_id string) ([]Todo, error) {
	query := `SELECT * from todo WHERE ticket_id=?`
	rows, err := adb.db.Queryx(query, ticket_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		value := Todo{}
		err = rows.StructScan(&value)
		if err != nil {
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
	return nil
}
