package main

import (
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	DB *DB
}

func (h *Handler) SayHello() string {
	return "Hello, world\n"
}

func (h *Handler) FindMessageById(params martini.Params) (int, string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return errorResponse(400, errors.New("You must submit an ID to lookup."))
	}

	message, err := h.DB.Find(id)

	if err == nil {
		return singleMessageResponse(message)
	} else {
		fmt.Println("Encountered an error fetching msg id=" + string(id) + ":")
		log.Fatal(err)
		return internalErrorResponse(err)
	}
}

func (h *Handler) FetchLatestMessages(req *http.Request) (int, string) {
	limit := req.URL.Query().Get("limit")
	if limit == "" { // no limit
		limit = "10"
	}
	fmt.Println("Fetching latest", limit, "messages")

	messages, err := h.DB.LatestMessages(limit)

	if err == nil {
		return messagesResponse(limit, messages)
	} else {
		fmt.Println("Encountered an error fetching the latest msgs:")
		log.Fatal(err)
		return internalErrorResponse(err)
	}
}

func (h *Handler) StoreMessage(msg Message) (int, string) {
	fmt.Println("Storing the following message:", msg.String())

	message, err := h.DB.InsertMessage(msg)

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
