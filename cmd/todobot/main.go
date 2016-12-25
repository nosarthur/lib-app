package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nosarthur/lib-app/server"
)

func main() {
	app := server.NewApplication()

	fmt.Println("start listening...")
	router := mux.NewRouter().StrictSlash(true)

	ticket := router.PathPrefix("/ticket").Subrouter()
	ticket.Handle("/add", server.AppHandler(app.AddTicket)).Methods("POST")
	ticket.Handle("/end/{id}", server.AppHandler(app.EndTicket)).Methods("DELETE")

	todo := router.PathPrefix("/todo").Subrouter()
	todo.Handle("/add", server.AppHandler(app.AddTodo)).Methods("POST")
	todo.Handle("/end/{ticket_id}/{idx}", server.AppHandler(app.EndTodo)).Methods("DELETE")

	router.Handle("/data", server.AppHandler(app.Get)).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
