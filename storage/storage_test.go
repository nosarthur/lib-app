package storage

import (
	"testing"
)

/*
	Database testing is performed on local sqlite database,
	instead of postgres database as in the production.
*/

func TestStorage(t *testing.T) {
	setupDB()
	t.Run("createTicket", func(t *testing.T) {})
	t.Run("readTicket", func(t *testing.T) {})
	t.Run("updateTicket", func(t *testing.T) {})
	t.Run("deleteTicket", func(t *testing.T) {})
	t.Run("createTodo", func(t *testing.T) {})
	t.Run("readTodo", func(t *testing.T) {})
	t.Run("updateTodo", func(t *testing.T) {})
	t.Run("deleteTodo", func(t *testing.T) {})
	t.Run("getTodoCount", func(t *testing.T) {})
	teardownDB()
}

func teardownDB() {
}

func setupDB() {
	var db AppDB
	db.MustCreateTables("sqlite3", "./test.sqlite")
}
