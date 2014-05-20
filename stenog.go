package main

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"net/http"
)

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Get("/api/messages/latest", fetchLatestMessages)
	m.Post("/api/messages/log", binding.Bind(Message{}), storeMessage)
	m.Run()
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

	messages, err := newDB().LatestMessages(limit)

	if err == nil {
		return messagesResponseMessage(limit, messages)
	} else {
		fmt.Println(err)
		return errorMessage(err)
	}
}

func storeMessage(msg Message) string {
	fmt.Println("Storing the following message:", msg.String())

	message, err := newDB().InsertMessage(msg)

	if err == nil {
		fmt.Println("Inserted message:", message)
		msgs := []Message{}
		msgs = append(msgs, message)
		return messagesResponseMessage("", msgs)
	} else {
		fmt.Println(err)
		return errorMessage(err)
	}
}
