package server

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	tickets, err := app.db.Get()
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
	var t storage.Ticket
	_ = json.NewDecoder(req.Body).Decode(&t)

	fmt.Println(req.Body)
	id, err := app.db.CreateTicket(&t)
	fmt.Println(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

func (app *Application) AddTodo(w http.ResponseWriter, req *http.Request) {
	var t storage.Todo
	_ = json.NewDecoder(req.Body).Decode(&t)

}

func (app *Application) EndTodo(w http.ResponseWriter, req *http.Request) {
}
func (app *Application) EndTicket(w http.ResponseWriter, req *http.Request) {
}
