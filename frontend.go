package gossip

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/parkr/gossip/database"
	"github.com/parkr/gossip/template"
)

func ensureLeadingHash(room string) string {
	if strings.HasPrefix(room, "#") {
		return room
	}
	return "#" + room
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	oldURL := r.URL
	r.URL = &url.URL{
		Scheme:     oldURL.Scheme,
		Opaque:     oldURL.Opaque,
		User:       oldURL.User,
		Host:       oldURL.Host,
		Path:       "/room/%23" + h.DefaultRoom(),
		RawPath:    "/room/%23" + h.DefaultRoom(),
		ForceQuery: oldURL.ForceQuery,
		RawQuery:   oldURL.RawQuery,
		Fragment:   oldURL.Fragment,
	}
	h.LatestMessagesByRoom(w, r)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "no search term given", http.StatusBadRequest)
		return
	}
	if len(query) < 3 {
		http.Error(w, "search term must be 3+ characters long", http.StatusBadRequest)
		return
	}

	cacheKey := "search-" + query
	messages, err := h.FetchAndCacheList(r, cacheKey, func() ([]database.Message, error) {
		return h.DB.ListByFuzzyMessage(query)
	})

	if err == sql.ErrNoRows || len(messages) == 0 {
		http.Error(w, "no results for "+query, http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "\n\ncouldn't fetch messages: %+v", err)
		http.Error(w, "couldn't fetch messages", http.StatusInternalServerError)
		return
	}

	messagesGroupedByRoom := map[string][]database.Message{}
	for _, message := range messages {
		if _, ok := messagesGroupedByRoom[message.Room]; !ok {
			messagesGroupedByRoom[message.Room] = []database.Message{}
		}
		messagesGroupedByRoom[message.Room] = append(messagesGroupedByRoom[message.Room], message)
	}
	data := &template.SearchTemplateData{
		Results: messagesGroupedByRoom,
		Total:   len(messages),
		Rooms:   h.AllRooms(),
		Query:   query,
	}
	if err := template.SearchTemplate.Execute(w, data); err != nil {
		fmt.Fprintf(w, "\n\n%+v", err)
	}
}

func (h *Handler) LatestMessagesByRoom(w http.ResponseWriter, r *http.Request) {
	unescapedURLPath, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		LogWithRequestID(r, fmt.Sprintf("Couldn't unescape URL.Path '%s': %+v", r.URL.Path, err))
		unescapedURLPath = r.URL.Path
	}
	room := ensureLeadingHash(strings.TrimPrefix(unescapedURLPath, "/room/"))
	limit := resultsLimit(r)

	cacheKey := "messages-by-room-" + room + "-" + fmt.Sprintf("%d", limit)
	messages, err := h.FetchAndCacheList(r, cacheKey, func() ([]database.Message, error) {
		return h.DB.LatestMessagesByRoom(room, limit)
	})

	if err == sql.ErrNoRows || len(messages) == 0 {
		http.Error(w, "no results for "+room, http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "\n\ncouldn't fetch messages: %+v", err)
		http.Error(w, "couldn't fetch messages", http.StatusInternalServerError)
		return
	}

	setLastModifiedAt(w, messages)
	data := &template.ListTemplateData{
		Messages:    messages,
		Rooms:       h.AllRooms(),
		CurrentRoom: room,
	}
	if err := template.ListTemplate.Execute(w, data); err != nil {
		fmt.Fprintf(w, "\n\n%+v", err)
	}
}

func (h *Handler) LatestMessagesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := strings.TrimPrefix(r.URL.Path, "/messages/by/")
	limit := resultsLimit(r)

	cacheKey := "messages-by-author-" + author + "-" + fmt.Sprintf("%d", limit)
	messages, err := h.FetchAndCacheList(r, cacheKey, func() ([]database.Message, error) {
		return h.DB.LatestMessagesByAuthor(author, limit)
	})

	if err == sql.ErrNoRows || len(messages) == 0 {
		http.Error(w, "no results for "+author, http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "\n\ncouldn't fetch messages: %+v", err)
		http.Error(w, "couldn't fetch messages", http.StatusInternalServerError)
		return
	}

	setLastModifiedAt(w, messages)
	data := &template.ListTemplateData{
		Messages:      messages,
		Rooms:         h.AllRooms(),
		CurrentAuthor: author,
	}
	if err := template.ListTemplate.Execute(w, data); err != nil {
		fmt.Fprintf(w, "\n\n%+v", err)
	}
}

func (h *Handler) MessageContext(w http.ResponseWriter, r *http.Request) {
	limit := 5
	messageIDStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/messages/"), "/context")
	if messageIDStr == "" {
		http.Error(w, "no message id given", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil || messageID == 0 {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	cacheKey := "message-" + messageIDStr
	message, err := h.FetchAndCacheGet(r, cacheKey, func() (*database.Message, error) {
		return h.DB.Find(messageID)
	})
	if err == sql.ErrNoRows {
		http.Error(w, "no message with id "+messageIDStr, http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "\n\ncouldn't fetch message: %+v", err)
		http.Error(w, "couldn't fetch message", http.StatusInternalServerError)
		return
	}

	priorCacheKey := "prior-" + messageIDStr
	priorMessages, err := h.FetchAndCacheList(r, priorCacheKey, func() ([]database.Message, error) {
		return h.DB.PriorMessages(message.Room, message.At, limit)
	})
	if err != nil && err != sql.ErrNoRows {
		fmt.Fprintf(w, "\n\ncouldn't fetch prior messages: %+v", err)
		http.Error(w, "couldn't fetch prior messages", http.StatusInternalServerError)
		return
	}

	subsequentCacheKey := "subsequent-" + messageIDStr
	subsequentMessages, err := h.FetchAndCacheList(r, subsequentCacheKey, func() ([]database.Message, error) {
		return h.DB.SubsequentMessages(message.Room, message.At, limit)
	})
	if err != nil && err != sql.ErrNoRows {
		fmt.Fprintf(w, "\n\ncouldn't fetch subsequent messages: %+v", err)
		http.Error(w, "couldn't fetch subsequent messages", http.StatusInternalServerError)
		return
	}

	setLastModifiedAt(w, subsequentMessages)
	data := &template.ShowTemplateData{
		PriorMessages:      priorMessages,
		Message:            *message,
		SubsequentMessages: subsequentMessages,
		Rooms:              h.AllRooms(),
		CurrentRoom:        message.Room,
	}
	if err := template.ShowTemplate.Execute(w, data); err != nil {
		fmt.Fprintf(w, "\n\n%+v", err)
	}
}

// Pulls the limit query parameter from the request
// and returns default of 20 if blank or non-integer.
func resultsLimit(r *http.Request) int {
	values := r.URL.Query()["limit"]
	if len(values) == 0 {
		return 20
	}

	limitStr := values[0]
	if limitStr == "" {
		return 20
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		LogWithRequestID(r, fmt.Sprintf("Bad limit '%s': %+v", limitStr, err))
		return 20
	}

	// Cap limit at a reasonable number for the DB.
	if limit > 100 {
		return 100
	}

	return limit
}
