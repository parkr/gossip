package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

const (
	InsertionQuery    = "INSERT INTO messages (room, author, message, at, created_at, updated_at) VALUES (:room, :author, :message, :at, NOW(), NOW())"
	SelectLatestQuery = "SELECT * FROM messages ORDER BY at DESC LIMIT 0,%s"
	SelectMessageById = "SELECT * FROM messages WHERE id = %d"
)

type DB struct {
	Connection *sqlx.DB
}

func dbUrl() string {
	username := os.Getenv("GOSSIP_DB_USERNAME")
	password := os.Getenv("GOSSIP_DB_PASSWORD")
	dbname := os.Getenv("GOSSIP_DB_DBNAME")
	return username + ":" + password + "@/" + dbname
}

func newDB() *DB {
	db := &DB{}
	conn, err := sqlx.Connect("mysql", dbUrl())
	if err != nil {
		fmt.Println("CRAP: couldn't connect to the database.")
		fmt.Println(err)
	}
	db.Connection = conn
	return db
}

func (db *DB) Close() error {
	err := db.Connection.DB.Close()
	db.Connection = nil
	return err
}

func (db *DB) Find(id int) (Message, error) {
	msg := Message{}
	err := db.Connection.Get(&msg, fmt.Sprintf(SelectMessageById, id))
	return msg, err
}

func (db *DB) LatestMessages(limit string) ([]Message, error) {
	messages := []Message{}
	err := db.Connection.Select(&messages, fmt.Sprintf(SelectLatestQuery, limit))
	return messages, err
}

func (db *DB) InsertMessage(msg Message) (Message, error) {
	result, err := db.Connection.NamedExec(InsertionQuery, msg.ForInsertion())
	if err != nil {
		return msg, err
	}
	lastInsertId, _ := result.LastInsertId()
	return db.Find(int(lastInsertId))
}
