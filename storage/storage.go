package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AppDB struct {
	db *sqlx.DB
}

func (adb *AppDB) MustInit() {
	adb.db = sqlx.MustConnect("sqlite3", "./app.sqlite")
	fmt.Println("database connected.")
}

func (adb *AppDB) GetAll() ([]Ticket, error) {
	rows, err := adb.db.Queryx("SELECT * FROM ticket")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tickets := []Ticket{}
	for rows.Next() {
		i := Ticket{}
		err = rows.StructScan(&i)
		if err != nil {
			panic(err)
		}
		tickets = append(tickets, i)
	}

	return tickets, rows.Err()
}

func (adb *AppDB) MustCreateTables() {
	schema := `
	CREATE TABLE ticket (
		id 			TEXT PRIMARY KEY,
		description TEXT,
		start_time  DATETIME NOT NULL,
		end_time    DATETIME,
		priority    INTEGER NOT NULL
	);

	CREATE TABLE todo (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		ticket_id TEXT NOT NULL,
		idx       INTEGER NOT NULL,
		item      TEXT NOT NULL,
		done      INTEGER NOT NULL
	);`
	db := sqlx.MustConnect("sqlite3", "./app.sqlite")
	_, err := db.Exec(schema)
	if err != nil {
		panic(err)
	}
	adb.db = db
}
