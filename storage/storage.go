package storage

import "github.com/jmoiron/sqlx"

type AppDB struct {
	db *sqlx.DB
}

func (adb *AppDB) MustInit() {
	adb.db = sqlx.MustConnect("sqlite3", "./app.sqlite")
}

func (adb *AppDB) Get() ([]Ticket, error) {
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
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		label       TEXT NOT NULL,
		description TEXT,
		start_time  DATETIME NOT NULL,
		end_time    DATETIME,
		priority    INTEGER
	);

	CREATE TABLE todo (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		ticket_id INTEGER NOT NULL,
		item      INTEGER NOT NULL,
		active    INTEGER NOT NULL
	);`
	db := sqlx.MustConnect("sqlite3", "./app.sqlite")
	_, err := db.Exec(schema)
	if err != nil {
		panic(err)
	}
	adb.db = db
}
