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

	fmt.Println("now listening...")
	router := mux.NewRouter()
	router.HandleFunc("/", app.Hello).Methods("GET")
	router.HandleFunc("/get", app.Get).Methods("GET")
	router.HandleFunc("/todo/create", app.CreateTodo).Methods("POST")
	router.HandleFunc("/ticket/create", app.CreateTicket).Methods("POST")
	router.HandleFunc("/ticket/end", app.EndTicket).Methods("DELETE")
	router.HandleFunc("/todo/end", app.EndTodo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
