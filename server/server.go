package server

import (
	"encoding/json"
	"fmt"
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
	fmt.Fprintln(w, "hello, nos!")
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
	fmt.Println(tickets)
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
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		panic(err)
	}

	err = app.db.CreateTicket(t)
	if err != nil {
		panic(err)
	}
}

func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	t, err := app.db.ReadTicket(params["id"])
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
	err = app.db.UpdateTicket(t)
	if err != nil {
		panic(err)
	}
}

func (app *Application) AddTodo(w http.ResponseWriter, req *http.Request) {
	t := storage.Todo{}
	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	id, err := app.db.CreateTodo(t)
	fmt.Println(id)
	if err != nil {
		panic(err)
	}
}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	_, err := app.db.ReadTicket(params["ticket_id"])
	if err != nil {
		panic(err)
	}
	idx, err := strconv.ParseInt(params["idx"], 10, 64)
	if err != nil {
		panic(err)
	}
	t, err := app.db.ReadTodo(params["ticket_id"], idx)
	t.Done = true
	err = app.db.UpdateTodo(t)
	if err != nil {
		panic(err)
	}
}
