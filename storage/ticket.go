package storage

import (
	"fmt"
	"time"
)

// Ticket contains Todos
type Ticket struct {
	Id        string     `db:"id" json:"id"`
	Detail    string     `db:"detail" json:"detail"`
	Todos     []*Todo    `json:"todos"`
	StartTime time.Time  `db:"start_time" json:"start_time"`
	EndTime   *time.Time `db:"end_time" json:"end_time"`
	Priority  bool       `db:"priority" json:"priority"`
}

func (adb *AppDB) CreateTicket(t Ticket) error {
	errMsg := fmt.Sprintf("Cannot create Ticket=%v", t)
	_, err := adb.db.NamedExec("INSERT INTO ticket (id, detail, start_time, priority) VALUES (:id, :detail, :start_time, :priority)", &t)
	if err != nil {
		return fmt.Errorf("%v, %v", errMsg, err)
	}
	adb.LastUpdate = time.Now()
	return nil
}

func (adb *AppDB) ReadTicket(id string) (Ticket, error) {
	t := Ticket{}
	query := `SELECT * FROM ticket WHERE id=$1`
	err := adb.db.Get(&t, query, id)
	return t, err
}

func (adb *AppDB) UpdateTicket(t Ticket) error {
	_, err := adb.db.NamedExec(`UPDATE ticket SET detail=:detail, start_time=:start_time, end_time=:end_time, priority=:priority WHERE id=:id;`, &t)
	if err != nil {
		return err
	}
	adb.LastUpdate = time.Now()
	return nil
}

// DeleteTicket deletes the ticket and the associated todos
func (adb *AppDB) DeleteTicket(id string) error {
	_, err := adb.db.Exec(`DELETE FROM todo WHERE ticket_id=$1;`, id)
	if err != nil {
		return err
	}
	_, err = adb.db.Exec(`DELETE FROM ticket WHERE id=$1;`, id)
	if err != nil {
		return err
	}
	adb.LastUpdate = time.Now()
	return nil
}
