package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nosarthur/lib-app/ticket"
)

type Application struct {
	db ticket.AppDB
}

func NewApplication() *Application {
	var app Application
	app.db.MustInit()
	return &app
}

func (app *Application) Get(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(ticket.Ticket{})
}

func (app *Application) CreateTicket(w http.ResponseWriter, req *http.Request) {
	var t ticket.Ticket
	_ = json.NewDecoder(req.Body).Decode(&t)

	test := ticket.Ticket{}
	fmt.Println(test)
	id, err := app.db.AddTicket(&test)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

func (app *Application) CreateTodo(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(ticket.Ticket{})
}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) {
}
func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) {
}
func (app *Application) Hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, nos!")
}
