package gossip

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/parkr/gossip/database"
)

type Handler struct {
	DB    *database.DB
	Cache *cache.Cache

	allRooms       []string
	skippedAuthors []string
	defaultRoom    string
}

func NewHandler() *Handler {
	return &Handler{
		DB: database.New(),
		// Create a cache with a default expiration time of 5 minutes, and which
		// purges expired items every 10 minutes
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/_health" {
		h.HealthCheck(w, r)
		return
	}

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

func (h *Handler) DefaultRoom() string {
	if h.defaultRoom == "" {
		h.defaultRoom = os.Getenv("GOSSIP_DEFAULT_ROOM") // this should NOT have a hash/pound symbol
	}
	return h.defaultRoom
}

type messagesListFunc func() ([]database.Message, error)

func (h *Handler) FetchAndCacheList(r *http.Request, key string, f messagesListFunc) ([]database.Message, error) {
	var messages []database.Message
	var err error

	messagesInterface, found := h.Cache.Get(key)
	if found {
		LogWithRequestID(r, fmt.Sprintf("Pulling %s results from cache", key))
		messages = messagesInterface.([]database.Message)
	} else {
		LogWithRequestID(r, fmt.Sprintf("Pulling %s results from database, stuffing in cache", key))
		messages, err = f()
		if err != nil {
			return messages, err // don't set the cache if there's an error
		}
		h.Cache.Set(key, messages, cache.DefaultExpiration)
	}

	return messages, err
}

type messagesGetFunc func() (*database.Message, error)

func (h *Handler) FetchAndCacheGet(r *http.Request, key string, f messagesGetFunc) (*database.Message, error) {
	var message *database.Message
	var err error

	messageInterface, found := h.Cache.Get(key)
	if found {
		LogWithRequestID(r, fmt.Sprintf("Pulling %s result from cache", key))
		message = messageInterface.(*database.Message)
	} else {
		LogWithRequestID(r, fmt.Sprintf("Pulling %s result from database, stuffing in cache", key))
		message, err = f()
		if err != nil {
			return message, err // don't set the cache if there's an error
		}
		h.Cache.Set(key, message, cache.DefaultExpiration)
	}

	return message, err
}
