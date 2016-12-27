package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	Database testing is performed on local postgres database.
*/

var adb AppDB

func setupDB() {
	adb.MustCreateTables("postgresql://localhost?sslmode=disable")
}

func teardownDB() {
	adb.MustDropTables()
}

func TestStorage(t *testing.T) {
	setupDB()
	t.Run("CreateTicket", func(t *testing.T) {
		tkt1 := Ticket{Id: "1", StartTime: time.Now()}
		tkt2 := Ticket{Id: "2", StartTime: time.Now()}
		tkt3 := Ticket{Id: "3", StartTime: time.Now()}
		err := adb.CreateTicket(tkt1)
		assert.Nil(t, err)
		err = adb.CreateTicket(tkt2)
		assert.Nil(t, err)
		err = adb.CreateTicket(tkt3)
		assert.Nil(t, err)
		// create an exisiting ticket
		err = adb.CreateTicket(tkt1)
		assert.NotNil(t, err)
	})
	t.Run("ReadTicket", func(t *testing.T) {
		tkt, err := adb.ReadTicket("1")
		assert.Nil(t, err)
		assert.Equal(t, "1", tkt.Id)
		assert.Equal(t, false, tkt.Priority)
		// read a non-existing ticket
		tkt, err = adb.ReadTicket("30")
		assert.NotNil(t, err)
	})
	t.Run("UpdateTicket", func(t *testing.T) {
		tkt := Ticket{Id: "2", Priority: true}
		err := adb.UpdateTicket(tkt)
		assert.Nil(t, err)
		tkt2, err := adb.ReadTicket("2")
		assert.Nil(t, err)
		assert.Equal(t, true, tkt2.Priority)
		// updating a non-existing ticket is ok
		tkt.Id = "20"
		err = adb.UpdateTicket(tkt)
		assert.Nil(t, err)
	})
	t.Run("DeleteTicket", func(t *testing.T) {
		err := adb.DeleteTicket("3")
		assert.Nil(t, err)
		// deleting a non-existing ticket is ok
		err = adb.DeleteTicket("20")
		assert.Nil(t, err)
	})
	t.Run("CreateTodo", func(t *testing.T) {
		t1 := Todo{TicketId: "1", Item: "test1"}
		t2 := Todo{TicketId: "2"}
		t30 := Todo{TicketId: "30"}
		err := adb.CreateTodo(t1)
		assert.Nil(t, err)
		t1.Item = "test2"
		err = adb.CreateTodo(t1)
		assert.Nil(t, err)
		t1.Item = "test3"
		err = adb.CreateTodo(t1)
		assert.Nil(t, err)
		err = adb.CreateTodo(t2)
		assert.Nil(t, err)
		// create a Todo with invalid TicketId
		err = adb.CreateTodo(t30)
		assert.NotNil(t, err)
	})
	t.Run("ReadTodo", func(t *testing.T) {
		todo, err := adb.ReadTodo("1", 1)
		assert.Nil(t, err)
		assert.Equal(t, "test1", todo.Item)
		todo, err = adb.ReadTodo("1", 3)
		assert.Nil(t, err)
		assert.Equal(t, "test3", todo.Item)
		// read a non-exisiting Todo
		todo, err = adb.ReadTodo("1", 10)
		assert.NotNil(t, err)
		// read a Todo with invalid ticket_id
		todo, err = adb.ReadTodo("10", 1)
		assert.NotNil(t, err)
	})
	t.Run("ReadTodos", func(t *testing.T) {
		todos, err := adb.ReadTodos("1")
		assert.Nil(t, err)
		assert.Equal(t, 3, len(todos))
		// read a Todo with invalid ticket_id
		todos, err = adb.ReadTodos("10")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(todos))
	})
	t.Run("UpdateTodo", func(t *testing.T) {
		t1 := Todo{TicketId: "1", Idx: 1, Item: "new item"}
		err := adb.UpdateTodo(t1)
		assert.Nil(t, err)
		t2, err := adb.ReadTodo(t1.TicketId, t1.Idx)
		assert.Nil(t, err)
		assert.Equal(t, t1.Item, t2.Item)
		// update a non-existing ticket is ok
		t1.Idx = 20
		err = adb.UpdateTodo(t1)
		assert.Nil(t, err)
	})
	t.Run("DeleteTodo", func(t *testing.T) {
		t2 := Todo{TicketId: "2", Idx: 1}
		err := adb.DeleteTodo(t2)
		assert.Nil(t, err)
		// deleting a non-existing ticket is ok
		t2.Id = 100
		err = adb.DeleteTodo(t2)
		assert.Nil(t, err)
	})
	t.Run("getTodoCount", func(t *testing.T) {
		n, err := adb.getTodoCount("1")
		assert.Nil(t, err)
		assert.Equal(t, int64(3), n)
		n, err = adb.getTodoCount("2")
		assert.Nil(t, err)
		assert.Equal(t, int64(0), n)
		// access the Todo count of a non-existing ticket
		n, err = adb.getTodoCount("10")
		assert.Nil(t, err)
		assert.Equal(t, int64(0), n)
	})
	teardownDB()
}
