package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
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
	return "Hello, world\n"
}

func findMessageById(params martini.Params) (int, string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return errorResponse(400, errors.New("You must submit an ID to lookup."))
	}

	db := newDB()
	message, err := db.Find(id)
	db.Close()
	if err == nil {
		return singleMessageResponse(message)
	} else {
		fmt.Println("Encountered an error fetching msg id=" + string(id) + ":")
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

	db := newDB()
	messages, err := db.LatestMessages(limit)
	db.Close()

	if err == nil {
		return messagesResponse(limit, messages)
	} else {
		fmt.Println("Encountered an error fetching the latest msgs:")
		log.Fatal(err)
		return internalErrorResponse(err)
	}
}

func storeMessage(msg Message) (int, string) {
	fmt.Println("Storing the following message:", msg.String())

	db := newDB()
	message, err := db.InsertMessage(msg)
	db.Close()

	if err == nil {
		fmt.Println("Inserted message:", message)
		msgs := []Message{}
		msgs = append(msgs, message)
		return messagesResponse("", msgs)
	} else {
		log.Fatal(err)
		return internalErrorResponse(err)
	}
}
