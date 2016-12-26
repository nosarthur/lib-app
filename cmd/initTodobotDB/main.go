package main

import (
	"os"

	"github.com/nosarthur/todobot/storage"
)

func main() {
	var db storage.AppDB
	db.MustCreateTables("postgres", os.Getenv("DATABASE_URL"))
}
