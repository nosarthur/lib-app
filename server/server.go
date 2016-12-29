/*
	package server implements http handlers, router, and slack bot handler
*/
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nosarthur/todobot/storage"
)

// appHandler decorates the http handlers with logging
type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := fn(w, req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// authHandler decorates the http handlers with authentication and logging
func authHandler(f func(http.ResponseWriter, *http.Request) error) appHandler {
	authf := func(w http.ResponseWriter, req *http.Request) error {
		if req.Header.Get("Token") != os.Getenv("Token") {
			return fmt.Errorf("Authentication failed.")
		}
		return f(w, req)
	}
	return appHandler(authf)
}

type application struct {
	db storage.AppDB
}

// NewApplication creates a new application with the associated database initialized
func NewApplication(dbURL string) *application {
	app := application{db: storage.AppDB{URL: dbURL}}
	app.db.MustInit()
	return &app
}

// NewRouter sets up the routes for the web application
func NewRouter(app *application) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// these four handlers require authentication
	ticket := router.PathPrefix("/ticket").Subrouter()
	ticket.Handle("/add", authHandler(app.AddTicket)).Methods("POST")
	ticket.Handle("/end/{id}", authHandler(app.EndTicket)).Methods("DELETE")

	todo := router.PathPrefix("/todo").Subrouter()
	todo.Handle("/add", authHandler(app.AddTodo)).Methods("POST")
	todo.Handle("/end/{ticket_id}/{idx}", authHandler(app.EndTodo)).Methods("DELETE")

	// these two handlers do not require authentication
	router.Handle("/data", appHandler(app.Data)).Methods("GET")
	router.Handle("/slack", appHandler(app.Slack)).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	return router
}
