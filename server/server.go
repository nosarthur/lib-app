package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nosarthur/todobot/storage"
)

type Application struct {
	db storage.AppDB
}

func NewApplication() *Application {
	var app Application
	app.db.MustInit()
	return &app
}

type AppHandler func(http.ResponseWriter, *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := fn(w, req); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (app *Application) Get(w http.ResponseWriter, req *http.Request) error {
	tickets, err := app.db.GetAll()
	if err != nil {
		return err
	}
	for i, t := range tickets {
		todos, err := app.db.ReadTodos(t.Id)
		if err != nil {
			return err
		}
		tickets[i].Todos = make([]storage.Todo, len(todos))
		copy(tickets[i].Todos, todos)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reply := map[string]interface{}{
		"tickets": tickets,
	}
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		return err
	}
	return nil
}

func (app *Application) AddTicket(w http.ResponseWriter, req *http.Request) error {
	t := storage.Ticket{StartTime: time.Now()}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	if err := app.db.CreateTicket(t); err != nil {
		return err
	}
	return nil
}

func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	t, err := app.db.ReadTicket(vars["id"])
	if err != nil {
		return err
	}
	if t.EndTime != nil {
		return errors.New("It's alread ended.")
	}
	now := time.Now()
	if now.Before(t.StartTime) {
		return errors.New("End time ealier than start time")
	}
	t.EndTime = &now
	if err = app.db.UpdateTicket(t); err != nil {
		return err
	}
	return nil
}

func (app *Application) AddTodo(w http.ResponseWriter, req *http.Request) error {
	t := storage.Todo{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	if err := app.db.CreateTodo(t); err != nil {
		return err
	}
	return nil
}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	_, err := app.db.ReadTicket(vars["ticket_id"])
	if err != nil {
		return err
	}
	idx, err := strconv.ParseInt(vars["idx"], 10, 64)
	if err != nil {
		return err
	}
	t, err := app.db.ReadTodo(vars["ticket_id"], idx)
	if err != nil {
		return err
	}
	t.Done = true
	if err = app.db.UpdateTodo(t); err != nil {
		return err
	}
	return nil
}
