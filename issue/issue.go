package issue

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Issue struct {
	Id          int64      `db:"id" json:"id"`
	Label       string     `db:"label" json:"label"`
	Description string     `db:"description" json:"description"`
	StartTime   time.Time  `db:"start_time" json:"start_time"`
	EndTime     *time.Time `db:"end_time" json:"end_time"`
	Priority    bool       `db:"priority" json:"priority"`
}
type Todo struct {
	Id      int64 `db:"id" json:"id"`
	Item    int64 `db:"item" json:"item"`
	IssueId int64 `db:"issue_id" json:"issue_id"`
	Active  bool  `db:"active" json:"active"`
}

type AppDB struct {
	db *sqlx.DB
}

func (adb *AppDB) MustInit() {
	adb.db = sqlx.MustConnect("sqlite3", "./app.db")
}

func (adb *AppDB) AddIssue(i *Issue) (int64, error) {
	result, err := adb.db.NamedExec("INSERT INTO issue (id, label, description, start_time, end_time, priority) VALUES (:id, :label, :description, :start_time, :end_time, :priority)", i)
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

func (adb *AppDB) Get() ([]Issue, error) {
	rows, err := adb.db.Queryx("SELECT * FROM issue")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	issues := []Issue{}
	for rows.Next() {
		i := Issue{}
		err = rows.StructScan(&i)
		if err != nil {
			panic(err)
		}
		issues = append(issues, i)
	}

	return issues, rows.Err()
}

func (adb *AppDB) MustCreate() {
	schema := `
	CREATE TABLE issue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		label TEXT NOT NULL,
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		priority INTEGER
	);

	CREATE TABLE todo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		issue_id INTEGER NOT NULL,
		item INTEGER NOT NULL,
		active INTEGER NOT NULL
	);`
	db := sqlx.MustConnect("sqlite3", "./app.db")
	_, err := db.Exec(schema)
	if err != nil {
		panic(err)
	}
	adb.db = db
}
