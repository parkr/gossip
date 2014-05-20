package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"time"
)

const (
	OkJsonMessage     = "{ ok: true }"
	FailedJsonMessage = "{ ok: false }"
)

type Message struct {
	Room       string `form:"room" json:"room" binding:"required"`
	Author     string `form:"author" json:"author" binding:"required"`
	Message    string `form:"message" json:"message" binding:"required"`
	At         string `form:"time" json:"time" binding:"required"`
	unexported string `form:"-"` // skip binding of unexported fields
}

func (msg *Message) String() string {
	return "<" + msg.Room + " by " + msg.Author + " at " + msg.At + ": " + msg.Message + ">"
}

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Post("/api/messages/log", binding.Bind(Message{}), storeMessage)
	m.Run()
}

func newDb() *sql.DB {
	db, err := sql.Open("mysql", "root@/witness")
	if err != nil {
		fmt.Println("CRAP the db couldn't be connected to.")
	}
	return db
}

func sayHello() string {
	return "Hello, world"
}

func storeMessage(msg Message) string {
	fmt.Println("Storing the following message:", msg.String())
	db := newDb()

	stmt, err := db.Prepare("INSERT INTO messages (room, author, message, at) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return FailedJsonMessage
	}
	defer stmt.Close() // in reality, you should check this call for error

	time, err := time.Parse("2006-01-02 15:04:05 -0700", msg.At)
	if err != nil {
		fmt.Println(err)
		return FailedJsonMessage
	}

	res, err := stmt.Exec(html.EscapeString(msg.Room), html.EscapeString(msg.Author), html.EscapeString(msg.Message), time)

	if err == nil {
		fmt.Sprintln(res)
		return OkJsonMessage
	} else {
		fmt.Println(err)
		return FailedJsonMessage
	}
}
