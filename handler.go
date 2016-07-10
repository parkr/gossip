package main

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/parkr/gossip/database"
	"github.com/parkr/gossip/response"
	"github.com/parkr/gossip/serializer"

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

func (h *Handler) FindMessagesByRoom(c web.C, w http.ResponseWriter, r *http.Request) {
	room := c.URLParams["room"]
	if room == "" {
		http.Error(w, "No room present.", 400)
		return
	}

	messages, err := h.DB.FindByRoom(room)

	if err != nil {
		errMsg := fmt.Sprintf("Could not fetch message id=%d: %s", id, err.Error())
		http.Error(w, response.New().WithError(errors.New(errMsg)).Json(), 500)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}

func (h *Handler) FindMessagesByAuthor(c web.C, w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) FindMessageWithContext(c web.C, w http.ResponseWriter, r *http.Request) {

}
