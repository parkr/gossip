package main

import (
	"fmt"
	"html"
	"time"
)

const (
	HowImUsedToItBeing   = "2006-01-02 15:04:05 -0700"
	HowJavaScriptSendsIt = "Mon, 02 Jan 2006 15:04:05 MST"
)

type Message struct {
	Id         int    `json:"id" db:"id"`
	Room       string `form:"room" json:"room" db:"room" binding:"required"`
	Author     string `form:"author" json:"author" db:"author" binding:"required"`
	Message    string `form:"text" json:"message" db:"message" binding:"required"`
	At         string `form:"time" json:"time" db:"at" binding:"required"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
	unexported string `form:"-"` // skip binding of unexported fields
}

type Messages []*Message

func (msg *Message) String() string {
	return "<" + msg.Room + " by " + msg.Author + " at " + msg.At + ": " + msg.Message + ">"
}

func (msg *Message) ForInsertion() map[string]interface{} {
	time, err := time.Parse(HowJavaScriptSendsIt, msg.At)
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"room":    html.EscapeString(msg.Room),
		"author":  html.EscapeString(msg.Author),
		"message": html.EscapeString(msg.Message),
		"at":      time,
	}
}
