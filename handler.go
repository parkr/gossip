package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/parkr/gossip/database"
	"github.com/parkr/gossip/response"
	"github.com/parkr/gossip/serializer"
)

var handler *Handler

type Handler struct {
	DB *database.DB
}

func init() {
	handler = &Handler{DB: database.New()}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		h.SayHello(w, r)
		return
	}

	if err := authenticate(r); err != nil {
		http.Error(w, response.New().WithError(err).Json(), http.StatusUnauthorized)
		return
	}

	switch r.URL.Path {
	case "/api/messages/latest":
		handler.FetchLatestMessages(w, r)
	case "/api/messages/log":
		handler.StoreMessage(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/api/messages/") {
			handler.FindMessageById(w, r)
		} else {
			http.Error(w, fmt.Sprintf("404 Not Found: %s", r.URL.Path), http.StatusNotFound)
		}
	}
}

func (h *Handler) SayHello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, there.\n")
}

func (h *Handler) FindMessageById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/messages/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "You must submit a numerical ID to lookup.", 400)
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

func (h *Handler) FetchLatestMessages(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	if limit == "" { // no limit
		limit = "10"
	}

	logForReq(r, fmt.Sprintf("Fetching latest %s messages", limit))

	messages, err := h.DB.LatestMessages(limit)

	if err != nil {
		errMsg := fmt.Sprintf("Could not fetch latest messages with limit=%s: %s", limit, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessages(messages).WithLimit(limit).Json())
}

func (h *Handler) StoreMessage(w http.ResponseWriter, r *http.Request) {
	var msg map[string]interface{}

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			fmt.Fprintf(w, response.New().WithError(err).Json())
			return
		}
	} else {
		room := r.PostFormValue("room")
		if room == "" {
			fmt.Fprintf(w, response.New().WithError(errors.New("No room specified. Skipping.")).Json())
			return
		}

		msg = map[string]interface{}{
			"room":    html.EscapeString(r.PostFormValue("room")),
			"author":  html.EscapeString(r.PostFormValue("author")),
			"message": html.EscapeString(r.PostFormValue("message")),
			"at":      serializer.ParseJavaScriptTime(r.PostFormValue("time")),
		}
	}

	logForReq(r, fmt.Sprintf("Inserting %+v", msg))

	message, err := h.DB.InsertMessage(msg)

	if err != nil {
		errMsg := fmt.Sprintf("Could not insert message %s: %s", msg, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}
