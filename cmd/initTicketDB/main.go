package main

import (
	"fmt"

	"github.com/nosarthur/lib-app/ticket"
)

func main() {
	var db ticket.AppDB
	db.MustCreateTables()
	fmt.Println("Tables created.")
}
