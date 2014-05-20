package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	InsertionQuery    = "INSERT INTO messages (room, author, message, at) VALUES (:room, :author, :message, :at)"
	SelectLatestQuery = "SELECT * FROM messages ORDER BY at DESC LIMIT 0,%s"
)

type DB struct {
	Connection *sqlx.DB
}

func newDB() *DB {
	db := &DB{}
	conn, err := sqlx.Connect("mysql", "root@/witness")
	if err != nil {
		fmt.Println("CRAP: couldn't connect to the database.")
		fmt.Println(err)
	}
	db.Connection = conn
	return db
}

func (db *DB) LatestMessages(limit string) ([]Message, error) {
	messages := []Message{}
	err := db.Connection.Select(&messages, fmt.Sprintf(SelectLatestQuery, limit))
	return messages, err
}

func (db *DB) InsertMessage(msg Message) (Message, error) {
	result, err := db.Connection.NamedExec(InsertionQuery, msg.ForInsertion())
	lastInsertId, _ := result.LastInsertId()
	msg.Id = int(lastInsertId)
	return msg, err
}
