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
	router.HandleFunc("/", app.Get).Methods("GET")
	router.HandleFunc("/todo/add", app.AddTodo).Methods("POST")
	router.HandleFunc("/todo/end", app.EndTodo).Methods("DELETE")
	router.HandleFunc("/ticket/add", app.AddTicket).Methods("POST")
	router.HandleFunc("/ticket/end", app.EndTicket).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
