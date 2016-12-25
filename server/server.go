package server

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	tickets, err := app.db.GetAll()
	if err != nil {
		panic(err)
	}
	for i, t := range tickets {
		todos, err := app.db.ReadTodos(t.Id)
		if err != nil {
			panic(err)
		}
		tickets[i].Todos = make([]storage.Todo, len(todos))
		copy(tickets[i].Todos, todos)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	reply := map[string]interface{}{
		"tickets": tickets,
	}
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}

func (app *Application) AddTicket(w http.ResponseWriter, req *http.Request) {
	t := storage.Ticket{StartTime: time.Now()}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if err := app.db.CreateTicket(t); err != nil {
		panic(err)
	}
}

func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	t, err := app.db.ReadTicket(vars["id"])
	if err != nil {
		panic(err)
	}
	if t.EndTime != nil {
		panic("It's alread ended.")
	}
	now := time.Now()
	if now.Before(t.StartTime) {
		panic("End time ealier than start time")
	}
	t.EndTime = &now
	if err = app.db.UpdateTicket(t); err != nil {
		panic(err)
	}
}

func (app *Application) AddTodo(w http.ResponseWriter, req *http.Request) {
	t := storage.Todo{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if err := app.db.CreateTodo(t); err != nil {
		panic(err)
	}
}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := app.db.ReadTicket(vars["ticket_id"])
	if err != nil {
		panic(err)
	}
	idx, err := strconv.ParseInt(vars["idx"], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
		panic(err)
	}
	t, err := app.db.ReadTodo(vars["ticket_id"], idx)
	t.Done = true
	if err = app.db.UpdateTodo(t); err != nil {
		panic(err)
	}
}
