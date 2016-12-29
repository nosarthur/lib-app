/*
	package storage implements CRUD for data types Ticket and Todo
*/

package storage

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// AppDB contains information of the database
type AppDB struct {
	db         *sqlx.DB
	URL        string
	LastUpdate time.Time
}

// MustInit connects to the database
func (adb *AppDB) MustInit() {
	adb.db = sqlx.MustConnect("postgres", adb.URL)
	log.Println("database connected.")
}

// GetAll returns all tickets in the database
// the todos are not included
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

// MustCreateTables creates the ticket and todo tables
func (adb *AppDB) MustCreateTables() {
	schema := `
	CREATE TABLE ticket (
		id          VARCHAR(16) PRIMARY KEY,
		detail      VARCHAR(32),
		start_time  TIMESTAMP WITH TIME ZONE  NOT NULL,
		end_time    TIMESTAMP WITH TIME ZONE,
		priority    BOOLEAN NOT NULL);
	CREATE TABLE todo (
		id        SERIAL PRIMARY KEY,
		ticket_id VARCHAR(16) NOT NULL,
		idx       INTEGER NOT NULL,
		item      VARCHAR(32) NOT NULL,
		done      BOOLEAN NOT NULL);`
	adb.db.MustExec(schema)
	log.Println("Tables created.")
}

// MustDropTables drops the ticket and todo tables
func (adb *AppDB) MustDropTables() {
	adb.db.MustExec(`DROP TABLE IF EXISTS ticket;`)
	adb.db.MustExec(`DROP TABLE IF EXISTS todo;`)
	log.Println("Tables dropped.")
}
