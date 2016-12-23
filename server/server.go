package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nosarthur/lib-app/issue"
)

type Application struct {
	db issue.AppDB
}

func NewApplication() *Application {
	var app Application
	//app.db.MustCreate()
	app.db.MustInit()
	return &app
}

func (app *Application) Get(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(issue.Issue{})
}

func (app *Application) CreateIssue(w http.ResponseWriter, req *http.Request) {
	var iss issue.Issue
	_ = json.NewDecoder(req.Body).Decode(&iss)

	test := issue.Issue{}
	fmt.Println(test)
	id, err := app.db.AddIssue(&test)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

func (app *Application) CreateTodo(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(issue.Issue{})
}

func (app *Application) DeleteTodo(w http.ResponseWriter, req *http.Request) {
}
func (app *Application) DeleteIssue(w http.ResponseWriter, req *http.Request) {
}
func (app *Application) Hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, nos!")
}
