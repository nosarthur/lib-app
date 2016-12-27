package storage

import "fmt"

type Todo struct {
	Id       int64  `db:"id" json:"id"`
	TicketId string `db:"ticket_id" json:"ticket_id"`
	Idx      int64  `db:"idx" json:"idx"`
	Item     string `db:"item" json:"item"`
	Done     bool   `db:"done" json:"done"`
}

func (adb *AppDB) getTodoCount(ticket_id string) (int64, error) {
	var count int64
	err := adb.db.QueryRow("SELECT COUNT(*) FROM todo WHERE ticket_id=$1", ticket_id).Scan(&count)
	return count, err
}

// CreateTodo sets Todo.Id and Todo.Idx automatically
func (adb *AppDB) CreateTodo(t Todo) error {
	errMsg := fmt.Sprintf("Cannot create Todo %v", t)

	if _, err := adb.ReadTicket(t.TicketId); err != nil {
		return fmt.Errorf("%v, %v", errMsg, err.Error())
	}
	count, err := adb.getTodoCount(t.TicketId)
	if err != nil {
		return fmt.Errorf("%v, %v", errMsg, err.Error())
	}
	t.Idx = count + 1
	_, err = adb.db.NamedExec("INSERT INTO todo (idx, item, ticket_id, done) VALUES (:idx, :item, :ticket_id, :done)", &t)
	if err != nil {
		return fmt.Errorf("%v, %v", errMsg, err.Error())
	}
	return nil
}

func (adb *AppDB) ReadTodo(ticket_id string, idx int64) (Todo, error) {
	errMsg := fmt.Sprintf("Cannot read Todo with ticket_id=%v and idx=%v", ticket_id, idx)
	t := Todo{}
	query := `SELECT * from todo WHERE ticket_id=$1 AND idx=$2`
	err := adb.db.Get(&t, query, ticket_id, idx)
	if err != nil {
		err = fmt.Errorf("%v, %v", errMsg, err.Error())
	}
	return t, err
}

func (adb *AppDB) ReadTodos(ticket_id string) ([]Todo, error) {
	errMsg := fmt.Sprintf("Cannot read Todos with ticket_id=%v", ticket_id)

	query := `SELECT * from todo WHERE ticket_id=$1`
	rows, err := adb.db.Queryx(query, ticket_id)
	if err != nil {
		return nil, fmt.Errorf("%v, %v", errMsg, err.Error())
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		value := Todo{}
		if err = rows.StructScan(&value); err != nil {
			return nil, fmt.Errorf("%v, %v", errMsg, err.Error())
		}
		todos = append(todos, value)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%v, %v", errMsg, rows.Err().Error())
	}
	return todos, nil
}

// UpdateTodo relies on the ticket_id and idx
func (adb *AppDB) UpdateTodo(t Todo) error {
	result, err := adb.db.NamedExec(`UPDATE todo SET item=:item, done=:done WHERE ticket_id=:ticket_id AND idx=:idx;`, &t)
	if err != nil {
		return err
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if nRows == 0 {
		return fmt.Errorf("Todo=%s does not exist for update", t)
	}
	return nil
}

// DeleteTodo relies on the ticket_id and idx
func (adb *AppDB) DeleteTodo(t Todo) error {
	_, err := adb.db.NamedExec(`DELETE FROM todo WHERE ticket_id=:ticket_id AND idx=:idx;`, &t)
	return err
}
