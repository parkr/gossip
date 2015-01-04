package database

import (
	"fmt"
)

type Message struct {
	Id        int    `json:"id" db:"id"`
	Room      string `json:"room" db:"room"`
	Author    string `json:"author" db:"author"`
	Message   string `json:"message" db:"message"`
	At        string `json:"time" db:"at"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func (msg *Message) String() string {
	return fmt.Sprintf("<%s by %s at %s: %s>", msg.Room, msg.Author, msg.At, msg.Message)
}
