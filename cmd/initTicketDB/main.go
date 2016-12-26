package main

import (
	"fmt"

	"github.com/nosarthur/todo-app/storage"
)

func main() {
	var db storage.AppDB
	db.MustCreateTables()
	fmt.Println("Tables created.")
}
