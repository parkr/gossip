package main

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
)

const (
	InsertionQuery = "INSERT INTO messages (room, author, message, at) VALUES (:room, :author, :message, :at)"
)

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Get("/api/messages/latest", fetchLatestMessages)
	m.Post("/api/messages/log", binding.Bind(Message{}), storeMessage)
	m.Run()
}

func fancyDb() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root@/witness")
	if err != nil {
		fmt.Println("CRAP the db couldn't be connected to.")
	}
	return db
}

func sayHello() string {
	return "Hello, world"
}

func fetchLatestMessages(req *http.Request) string {
	limit := req.URL.Query().Get("limit")
	if limit == "" { // no limit
		limit = "10"
	}
	fmt.Println("Fetching latest", limit, "messages")

	db := fancyDb()

	messages := []Message{}
	err := db.Select(&messages, "SELECT * FROM messages ORDER BY at DESC LIMIT 0,"+limit)

	if err == nil {
		return messagesResponseMessage(limit, messages)
	} else {
		fmt.Println(err)
		return errorMessage(err)
	}
}

func storeMessage(msg Message) string {
	fmt.Println("Storing the following message:", msg.String())
	db := fancyDb()

	ah, err := db.NamedExec(InsertionQuery, msg.ForInsertion())

	if err == nil {
		fmt.Sprintln(ah)
		return basicResponseMessage("true")
	} else {
		fmt.Println(err)
		return errorMessage(err)
	}
}
