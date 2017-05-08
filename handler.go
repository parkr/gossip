package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/parkr/gossip/database"
)

var handler *Handler

type Handler struct {
	DB *database.DB

	// cache!
	allRooms []string
}

func init() {
	handler = &Handler{DB: database.New()}
}

func (h *Handler) AllRooms() []string {
	if h.allRooms == nil {
		var err error
		h.allRooms, err = h.DB.AllRooms()
		if err != nil {
			log.Printf("error fetching rooms: %+v", err)
		}
	}

	return h.allRooms
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		h.Index(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/") {
		h.API(w, r)
		return
	}

	if r.URL.Path == "/search" {
		h.Search(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/room/") {
		h.LatestMessagesByRoom(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/messages/by/") {
		h.LatestMessagesByAuthor(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/messages/") && strings.HasSuffix(r.URL.Path, "/context") {
		h.MessageContext(w, r)
		return
	}

	http.Error(w, fmt.Sprintf("404 Not Found: %s", r.URL.Path), http.StatusNotFound)
}
