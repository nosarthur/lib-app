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

// logHandler decorates the http handlers with logging
type logHandler func(http.ResponseWriter, *http.Request) error

func (fn logHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := fn(w, req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// authHandler decorates the logHandlers with authentication
func authHandler(fn logHandler) logHandler {
	authf := func(w http.ResponseWriter, req *http.Request) error {
		if req.Header.Get("Token") != os.Getenv("Token") {
			return fmt.Errorf("Authentication failed.")
		}
		return fn(w, req)
	}
	return logHandler(authf)
}

type application struct {
	db          storage.AppDB
	slackRoutes map[string]logHandler
}

// NewApplication creates a new application with the associated database initialized
func NewApplication(dbURL string) *application {
	app := application{db: storage.AppDB{URL: dbURL}}
	app.db.MustInit()
	app.slackRoutes = map[string]logHandler{
		"/ticket/add": authHandler(logHandler(app.AddTicket)),
		"/ticket/end": authHandler(logHandler(app.EndTicket)),
		"/todo/add":   authHandler(logHandler(app.AddTodo)),
		"/todo/end":   authHandler(logHandler(app.EndTodo)),
	}
	return &app
}

// NewRouter sets up the routes for the web application
func NewRouter(app *application) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for url, handler := range app.slackRoutes {
		router.Handle(url, handler).Methods("POST")
	}
	//router.Handle("/ticket/end/{id}", ).Methods("POST")
	//	todo.Handle("/end/{ticket_id}/{idx}", ).Methods("DELETE")

	router.Handle("/data", logHandler(app.Data)).Methods("GET")
	router.Handle("/slack", logHandler(app.Slack)).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	return router
}
