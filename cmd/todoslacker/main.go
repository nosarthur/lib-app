package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nosarthur/todoslacker/server"
)

func main() {
	app := server.NewApplication(os.Getenv("DATABASE_URL"))
	router := server.NewRouter(app)

	log.Println("start listening...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
