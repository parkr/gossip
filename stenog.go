package main

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

func main() {
	// Setup Martini
	m := martini.Classic()
	m.Get("/", sayHello)
	m.Group("/api/messages", func(r martini.Router) {
		r.Get("/(?P<id>[0-9]+)", TokenAuthHandler(), findMessageById)
		r.Get("/latest", TokenAuthHandler(), fetchLatestMessages)
		r.Post("/log", TokenAuthHandler(), binding.Bind(Message{}), storeMessage)
	})
	m.Run()
}

func sayHello() string {
	return "Hello, world"
}

func findMessageById(params martini.Params) (int, string) {
	id := params["id"]
	if id == nil {
		return errorResponse(400, "You must submit an ID to lookup.")
	}
	message, err := newDB().Find(id)
	if err == nil {
		return singleMessageResponse(message)
	} else {
		fmt.Println("Encountered an error fetching msg id=" + id + ":")
		log.Fatal(err)
		return internalErrorResponse(err)
	}
}

func fetchLatestMessages(req *http.Request) (int, string) {
	limit := req.URL.Query().Get("limit")
	if limit == "" { // no limit
		limit = "10"
	}
	fmt.Println("Fetching latest", limit, "messages")

	messages, err := newDB().LatestMessages(limit)

	if err == nil {
		return 200, messagesResponseMessage(limit, messages)
	} else {
		fmt.Println("Encountered an error fetching the latest msgs:")
		log.Fatal(err)
		return errorMessage(err)
	}
}

func storeMessage(msg Message) (int, string) {
	fmt.Println("Storing the following message:", msg.String())

	message, err := newDB().InsertMessage(msg)

	if err == nil {
		fmt.Println("Inserted message:", message)
		msgs := []Message{}
		msgs = append(msgs, message)
		return 200, messagesResponseMessage("", msgs)
	} else {
		log.Fatal(err)
		return errorMessage(err)
	}
}
