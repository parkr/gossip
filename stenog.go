package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"html"
	"net/http"
	"time"
)

const (
	OkJsonMessage     = "{ ok: true }"
	FailedJsonMessage = "{ ok: false }"
)

type Message struct {
	Id         int    `json:"id" db:"id"`
	Room       string `form:"room" json:"room" db:"room" binding:"required"`
	Author     string `form:"author" json:"author" db:"author" binding:"required"`
	Message    string `form:"message" json:"message" db:"message" binding:"required"`
	At         string `form:"time" json:"time" binding:"required"`
	unexported string `form:"-"` // skip binding of unexported fields
}

type Messages []*Message

func (msg *Message) String() string {
	return "<" + msg.Room + " by " + msg.Author + " at " + msg.At + ": " + msg.Message + ">"
}

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Get("/messages/latest", fetchLatestMessages)
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

func fancyDb() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "dbname=bar")
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
	fmt.Println("limit =", limit)

	db := fancyDb()

	messages := []Message{}
	err := db.Select(&messages, "SELECT * FROM messages ORDER BY at DESC LIMIT 0,"+limit)
	fmt.Println(messages)

	if err != nil {
		json, _ := json.Marshal(messages)

		return "{ok:true,limit:" + limit + ",values:" + string(json) + "}"
	} else {
		fmt.Println(messages)
		fmt.Println(err)
		return FailedJsonMessage
	}
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
