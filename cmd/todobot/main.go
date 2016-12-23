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
	router.HandleFunc("/", app.Get).Methods("GET")

	ticket := router.PathPrefix("/ticket").Subrouter()
	ticket.HandleFunc("/add", app.AddTicket).Methods("POST")
	ticket.HandleFunc("/end/{id}", app.EndTicket).Methods("DELETE")

	todo := router.PathPrefix("/todo").Subrouter()
	todo.HandleFunc("/add", app.AddTodo).Methods("POST")
	todo.HandleFunc("/end", app.EndTodo).Methods("DELETE")
	//http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
