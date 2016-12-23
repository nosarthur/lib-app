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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
