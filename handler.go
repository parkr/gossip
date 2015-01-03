package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gossip/serializer"

	"github.com/zenazn/goji/web"
)

type Handler struct {
	DB *DB
}

func (h *Handler) SayHello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, there.\n")
}

func (h *Handler) FindMessageById(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		http.Error(w, "You must submit an ID to lookup.", 400)
		return
	}

	message, err := h.DB.Find(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch message id=%d: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, serializer.MarshalJson(map[string]interface{}{
		"messages": message,
	}))
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
