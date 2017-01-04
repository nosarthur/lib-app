package main

import (
	"os"

	"github.com/nosarthur/todoslacker/storage"
)

func main() {
	db := storage.AppDB{URL: os.Getenv("DATABASE_URL")}
	db.MustInit()
	db.MustCreateTables()
}
