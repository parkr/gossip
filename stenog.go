package main

import (
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
	InsertionQuery    = "INSERT INTO messages (room, author, message, at) VALUES (:room, :author, :message, :at)"
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

func (msg *Message) ForInsertion() map[string]interface{} {
	time, err := time.Parse("2006-01-02 15:04:05 -0700", msg.At)
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"room":    html.EscapeString(msg.Room),
		"author":  html.EscapeString(msg.Author),
		"message": html.EscapeString(msg.Message),
		"at":      time,
	}
}

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Get("/messages/latest", fetchLatestMessages)
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
	db := fancyDb()

	ah, err := db.NamedExec(InsertionQuery, msg.ForInsertion())

	if err == nil {
		fmt.Sprintln(ah)
		return OkJsonMessage
	} else {
		fmt.Println(err)
		return FailedJsonMessage
	}
}
