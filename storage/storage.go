package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AppDB struct {
	db *sqlx.DB
}

func (adb *AppDB) MustInit(dbLoc string) {
	adb.db = sqlx.MustConnect("postgres", dbLoc)
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

func (adb *AppDB) MustCreateTables(dbLoc string) {
	schema := `
	CREATE TABLE ticket (
		id          VARCHAR(16) PRIMARY KEY,
		detail      VARCHAR(32),
		start_time  TIMESTAMP WITH TIME ZONE  NOT NULL,
		end_time    TIMESTAMP WITH TIME ZONE,
		priority    BOOLEAN NOT NULL
	);

	CREATE TABLE todo (
		id        SERIAL PRIMARY KEY,
		ticket_id VARCHAR(16) NOT NULL,
		idx       INTEGER NOT NULL,
		item      VARCHAR(32) NOT NULL,
		done      BOOLEAN NOT NULL
	);`
	db := sqlx.MustConnect("postgres", dbLoc)
	db.MustExec(schema)
	adb.db = db
	fmt.Println("Tables created.")
}

func (adb *AppDB) MustDropTables() {
	adb.db.MustExec(`DROP TABLE ticket; DROP TABLE todo;`)
	fmt.Println("Tables dropped.")
}
