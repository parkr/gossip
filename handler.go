package main

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"gossip/database"
	"gossip/response"
	"gossip/serializer"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func logForReq(c web.C, message string) {
	log.Printf("[%s] %s", middleware.GetReqID(c), message)
}

var handler *Handler

type Handler struct {
	DB *database.DB
}

func init() {
	handler = &Handler{DB: database.New()}
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
		errMsg := fmt.Sprintf("Could not fetch message id=%d: %s", id, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}

func (h *Handler) FetchLatestMessages(c web.C, w http.ResponseWriter, r *http.Request) {
	limit := c.URLParams["limit"]
	if limit == "" { // no limit
		limit = "10"
	}

	logForReq(c, fmt.Sprintf("Fetching latest %s messages", limit))

	messages, err := h.DB.LatestMessages(limit)

	if err != nil {
		errMsg := fmt.Sprintf("Could not fetch latest messages with limit=%s: %s", limit, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessages(messages).WithLimit(limit).Json())
}

func (h *Handler) StoreMessage(c web.C, w http.ResponseWriter, r *http.Request) {
	room := r.PostFormValue("room")
	if room == "" {
		fmt.Fprintf(w, response.New().WithError(errors.New("No room specified. Skipping.")).Json())
		return
	}

	msg := map[string]interface{}{
		"room":    html.EscapeString(r.PostFormValue("room")),
		"author":  html.EscapeString(r.PostFormValue("author")),
		"message": html.EscapeString(r.PostFormValue("message")),
		"at":      serializer.ParseJavaScriptTime(r.PostFormValue("time")),
	}

	logForReq(c, fmt.Sprintf("Inserting %s:", msg))

	message, err := h.DB.InsertMessage(msg)

	if err != nil {
		errMsg := fmt.Sprintf("Could not insert message %s: %s", msg, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}
