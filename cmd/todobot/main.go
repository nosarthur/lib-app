package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nosarthur/todobot/server"
)

func main() {
	app := server.NewApplication()
	router := server.NewRouter(app)

	log.Println("start listening...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
