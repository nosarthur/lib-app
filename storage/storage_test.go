package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	Database testing is performed on local postgres database,
*/

var adb AppDB

func TestStorage(t *testing.T) {
	setupDB()
	t.Run("CreateTicket", func(t *testing.T) {
		tkt1 := Ticket{Id: "1", StartTime: time.Now()}
		tkt2 := Ticket{Id: "2", StartTime: time.Now()}
		err := adb.CreateTicket(tkt1)
		assert.Nil(t, err)
		err = adb.CreateTicket(tkt2)
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
		tkt, err = adb.ReadTicket("2")
		assert.Nil(t, err)
		assert.Equal(t, true, tkt.Priority)
		// update a non-existing ticket
		tkt.Id = "20"
		err = adb.UpdateTicket(tkt)
		assert.NotNil(t, err)
	})
	t.Run("DeleteTicket", func(t *testing.T) {
		err := adb.DeleteTicket("2")
		assert.Nil(t, err)
		// deleting a non-existing ticket is fine
		err = adb.DeleteTicket("20")
		assert.Nil(t, err)
	})
	t.Run("CreateTodo", func(t *testing.T) {
	})
	t.Run("ReadTodo", func(t *testing.T) {
	})
	t.Run("UpdateTodo", func(t *testing.T) {
	})
	t.Run("DeleteTodo", func(t *testing.T) {
	})
	t.Run("GetTodoCount", func(t *testing.T) {
	})
	teardownDB()
}

func teardownDB() {
	adb.MustDropTables()
}

func setupDB() {
	adb.MustCreateTables("postgresql://localhost?sslmode=disable")
}
