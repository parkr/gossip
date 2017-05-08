package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/parkr/gossip/database"
)

var handler *Handler

type Handler struct {
	DB *database.DB

	allRooms       []string
	skippedAuthors []string
}

func init() {
	handler = &Handler{DB: database.New()}
}

func (h *Handler) AllRooms() []string {
	if h.allRooms == nil {
		if rooms := os.Getenv("GOSSIP_ROOMS"); rooms != "" {
			// pull rooms in from the env
			h.allRooms = strings.Split(rooms, ",")
		} else {
			// pull rooms in from the DB
			var err error
			h.allRooms, err = h.DB.AllRooms()
			if err != nil {
				log.Printf("error fetching rooms: %+v", err)
			}
		}
	}

	return h.allRooms
}

func (h *Handler) SkippedAuthors() []string {
	if h.skippedAuthors == nil {
		h.skippedAuthors = strings.Split(os.Getenv("GOSSIP_SKIPPED_AUTHORS"), ",")
	}
	return h.skippedAuthors
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
