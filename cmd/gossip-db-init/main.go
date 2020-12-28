package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/parkr/gossip/database"
)

func main() {
	var retries int
	flag.IntVar(&retries, "retries", 10, "Number of times to retry the connection")
	flag.Parse()

	db := database.New()
	var err error
	for i := 1; i <= retries; i++ {
		log.Printf("attempting to connect (%d/%d)...", i, retries)
		_, err = db.Connect(context.Background())
		if err != nil {
			log.Printf("unable to connect: %v", err)
			time.Sleep(1 * time.Second)
		} else {
			log.Println("connection established")
			break
		}
	}

	if db.Connection == nil {
		log.Fatal("unable to establish connection to db")
	}

	err = db.InitDB(context.Background())
	if err != nil {
		log.Fatalf("unable to initialize db table: %v", err)
	} else {
		log.Println("db initialized")
	}
}
