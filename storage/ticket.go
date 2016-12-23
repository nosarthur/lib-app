package storage

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Ticket struct {
	Id          int64      `db:"id" json:"id"`
	Label       string     `db:"label" json:"label"`
	Description string     `db:"description" json:"description"`
	StartTime   time.Time  `db:"start_time" json:"start_time"`
	EndTime     *time.Time `db:"end_time" json:"end_time"`
	Priority    bool       `db:"priority" json:"priority"`
}

func (adb *AppDB) CreateTicket(t *Ticket) (int64, error) {
	result, err := adb.db.NamedExec("INSERT INTO ticket (id, label, description, start_time, end_time, priority) VALUES (:id, :label, :description, :start_time, :end_time, :priority)", t)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
