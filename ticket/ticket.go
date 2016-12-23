package ticket

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
type Todo struct {
	Id       int64 `db:"id" json:"id"`
	Item     int64 `db:"item" json:"item"`
	TicketId int64 `db:"ticket_id" json:"ticket_id"`
	Active   bool  `db:"active" json:"active"`
}

type AppDB struct {
	db *sqlx.DB
}

func (adb *AppDB) MustInit() {
	adb.db = sqlx.MustConnect("sqlite3", "./app.sqlite")
}

func (adb *AppDB) AddTicket(i *Ticket) (int64, error) {
	result, err := adb.db.NamedExec("INSERT INTO ticket (id, label, description, start_time, end_time, priority) VALUES (:id, :label, :description, :start_time, :end_time, :priority)", i)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (adb *AppDB) AddTodo(t Todo) {
	fmt.Println(t)

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
