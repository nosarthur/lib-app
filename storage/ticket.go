package storage

import (
	"time"
)

type Ticket struct {
	Id        string     `db:"id" json:"id"`
	Detail    string     `db:"detail" json:"detail"`
	Todos     []Todo     `json:"todos"`
	StartTime time.Time  `db:"start_time" json:"start_time"`
	EndTime   *time.Time `db:"end_time" json:"end_time"`
	Priority  bool       `db:"priority" json:"priority"`
}

func (adb *AppDB) CreateTicket(t Ticket) error {
	_, err := adb.db.NamedExec("INSERT INTO ticket (id, detail, start_time, priority) VALUES (:id, :detail, :start_time, :priority)", &t)
	return err
}

func (adb *AppDB) ReadTicket(id string) (Ticket, error) {
	t := Ticket{}
	query := `SELECT * FROM ticket WHERE id=?`
	err := adb.db.Get(&t, query, id)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (adb *AppDB) UpdateTicket(t Ticket) error {
	_, err := adb.db.NamedExec(`UPDATE ticket SET detail=:detail, start_time=:start_time, end_time=:end_time, priority=:priority WHERE id=:id;`, &t)
	return err
}

func (adb *AppDB) DeleteTicket(t Ticket) error {
	return nil
}
