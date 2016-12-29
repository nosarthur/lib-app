package main

import (
	"os"

	"github.com/nosarthur/todobot/storage"
)

func main() {
	db := storage.AppDB{URL: os.Getenv("DATABASE_URL")}
	db.MustInit()
	db.MustCreateTables()
}
