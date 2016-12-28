package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nosarthur/todobot/storage"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := fn(w, req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type application struct {
	db storage.AppDB
}

func NewApplication() *application {
	var app application
	app.db.MustInit(os.Getenv("DATABASE_URL"))
	return &app
}

func NewRouter(app *application) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	ticket := router.PathPrefix("/ticket").Subrouter()
	ticket.Handle("/add", appHandler(app.AddTicket)).Methods("POST")
	ticket.Handle("/end/{id}", appHandler(app.EndTicket)).Methods("DELETE")

	todo := router.PathPrefix("/todo").Subrouter()
	todo.Handle("/add", appHandler(app.AddTodo)).Methods("POST")
	todo.Handle("/end/{ticket_id}/{idx}", appHandler(app.EndTodo)).Methods("DELETE")

	router.Handle("/data", appHandler(app.Get)).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	return router
}

/*
	REST API without authentication

	http:get:: /data

	Example response:

	{
		"tickets" :[ {"id":"grocery",
		              "detail":"vegi",
					  "todos": []
					  "start_time": 2016-12-27T05:30:34.645428Z",
					  "end_time":null,
					  "priority":false}
		]
	}

*/
func (app *application) Get(w http.ResponseWriter, req *http.Request) error {
	tickets, err := app.db.GetAll()
	if err != nil {
		return err
	}
	for i, t := range tickets {
		todos, err := app.db.ReadTodos(t.Id)
		if err != nil {
			return err
		}
		tickets[i].Todos = make([]*storage.Todo, len(todos))
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

// http:post::/ticket/add/
func (app *application) AddTicket(w http.ResponseWriter, req *http.Request) error {
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

// http:delete::/ticket/end/id
func (app *application) EndTicket(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	t, err := app.db.ReadTicket(vars["id"])
	errMsg := fmt.Sprintf("Cannot end Ticket=%v.", t)
	if err != nil {
		return fmt.Errorf("%v %v", errMsg, err)
	}
	if t.EndTime != nil {
		return fmt.Errorf("%v It has ended already.", errMsg)
	}
	now := time.Now()
	if now.Before(t.StartTime) {
		return fmt.Errorf("%v Causality broken.", errMsg)
	}
	t.EndTime = &now
	if err = app.db.UpdateTicket(t); err != nil {
		return fmt.Errorf("%v %v", errMsg, err)
	}
	return nil
}

// http:post::/todo/add
func (app *application) AddTodo(w http.ResponseWriter, req *http.Request) error {
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

// http:delete::/todo/ticket_id/idx
func (app *application) EndTodo(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	if _, err := app.db.ReadTicket(vars["ticket_id"]); err != nil {
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
