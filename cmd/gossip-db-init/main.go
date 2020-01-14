package main

import (
	"context"
	"flag"
	"log"

	"github.com/parkr/gossip/database"
)

func main() {
	flag.Parse()

	db := database.NewWithInit()
	err := db.InitDB(context.Background())
	if err != nil {
		log.Fatalf("unable to initialize db table: %#v", err)
	} else {
		log.Println("db initialized")
	}
}
