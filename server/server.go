package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nosarthur/lib-app/storage"
)

type Application struct {
	db storage.AppDB
}

func NewApplication() *Application {
	var app Application
	app.db.MustInit()
	return &app
}

func (app *Application) Get(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello, nos!")
	tickets, err := app.db.GetAll()
	if err != nil {
		panic(err)
	}
	reply := map[string]interface{}{
		"tickets": tickets,
	}
	data, err := json.Marshal(reply)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func (app *Application) AddTicket(w http.ResponseWriter, req *http.Request) {
	t := storage.Ticket{StartTime: time.Now()}
	_ = json.NewDecoder(req.Body).Decode(&t)

	err := app.db.CreateTicket(t)
	if err != nil {
		panic(err)
	}
}

func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	t, err := app.db.ReadTicket(params["id"])
	if t.EndTime != nil {
		panic("It's alread ended.")
	}
	now := time.Now()
	if now.Before(t.StartTime) {
		panic("End time ealier than start time")
	}
	t.EndTime = &now
	err = app.db.UpdateTicket(t)
	if err != nil {
		panic(err)
	}
}

func (app *Application) AddTodo(w http.ResponseWriter, req *http.Request) {
	var t storage.Todo
	_ = json.NewDecoder(req.Body).Decode(&t)

}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) {
}
