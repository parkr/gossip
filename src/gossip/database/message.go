package database

import (
	"gossip/serializer"

	"fmt"
	"html"
)

type Message struct {
	Id         int    `json:"id" db:"id"`
	Room       string `json:"room" db:"room"`
	Author     string `json:"author" db:"author"`
	Message    string `json:"message" db:"message"`
	At         string `json:"time" db:"at"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
	unexported string `form:"-"` // skip binding of unexported fields
}

type Messages []*Message

func (msg *Message) String() string {
	return fmt.Sprintf("<%s by %s at %s: %s>", msg.Room, msg.Author, msg.At, msg.Message)
}

func (msg *Message) ForInsertion() map[string]interface{} {
	return map[string]interface{}{
		"room":    html.EscapeString(msg.Room),
		"author":  html.EscapeString(msg.Author),
		"message": html.EscapeString(msg.Message),
		"at":      serializer.ParseJavaScriptTime(msg.At),
	}
}
