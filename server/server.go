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
	"github.com/nosarthur/todoslacker/storage"
)

// logHandler decorates the http handlers with logging
type logHandler func(http.ResponseWriter, *http.Request) error

func (fn logHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := fn(w, req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// auth decorates the logHandlers with authentication
func auth(fn logHandler) logHandler {
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
		"/ticket/add": auth(logHandler(app.AddTicket)),
		"/ticket/end": auth(logHandler(app.EndTicket)),
		"/todo/add":   auth(logHandler(app.AddTodo)),
		"/todo/end":   auth(logHandler(app.EndTodo)),
	}
	return &app
}

// NewRouter sets up the routes for the web application
func NewRouter(app *application) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for url, handler := range app.slackRoutes {
		router.Handle(url, handler).Methods("POST")
	}
	router.Handle("/data", logHandler(app.Data)).Methods("GET")
	router.Handle("/slack", logHandler(app.Slack)).Methods("POST")

	router.Handle("/todo/delete", auth(logHandler(app.DeleteTodo))).Methods("DELETE")
	router.Handle("/ticket/delete", auth(logHandler(app.DeleteTicket))).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	return router
}
