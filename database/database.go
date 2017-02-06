package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	InsertionQuery    = "INSERT INTO messages (room, author, message, at, created_at, updated_at) VALUES (:room, :author, :message, :at, NOW(), NOW())"
	SelectLatestQuery = "SELECT * FROM messages ORDER BY at DESC LIMIT 0,%s"
	SelectMessageById = "SELECT * FROM messages WHERE id = %d"
)

type DB struct {
	Connection *sqlx.DB
}

var cachedDatabaseURL string

func databaseURL() string {
	if cachedDatabaseURL == "" {
		username := os.Getenv("GOSSIP_DB_USERNAME")
		password := os.Getenv("GOSSIP_DB_PASSWORD")
		os.Setenv("GOSSIP_DB_PASSWORD", "") // Unset this variable so it doesn't leak.
		dbname := os.Getenv("GOSSIP_DB_DBNAME")
		hostname := os.Getenv("GOSSIP_DB_HOSTNAME")
		cachedDatabaseURL = username + ":" + password + "@" + hostname + "/" + dbname
	}

	return cachedDatabaseURL
}

func New() *DB {
	return &DB{}
}

func (db *DB) GetConnection() *sqlx.DB {
	if db.Connection == nil {
		db.Connection = sqlx.MustConnect("mysql", databaseURL())
	}

	return db.Connection
}

func (db *DB) Close() error {
	if db.Connection == nil {
		return nil
	}

	err := db.Connection.DB.Close()
	db.Connection = nil
	return err
}

func (db *DB) Find(id int) (*Message, error) {
	msg := &Message{}
	err := db.GetConnection().Get(msg, fmt.Sprintf(SelectMessageById, id))
	return msg, err
}

func (db *DB) LatestMessages(limit string) ([]Message, error) {
	messages := []Message{}
	err := db.GetConnection().Select(&messages, fmt.Sprintf(SelectLatestQuery, limit))
	return messages, err
}

func (db *DB) InsertMessage(message map[string]interface{}) (*Message, error) {
	result, err := db.GetConnection().NamedExec(InsertionQuery, message)
	if err != nil {
		return nil, err
	}
	lastInsertId, _ := result.LastInsertId()
	return db.Find(int(lastInsertId))
}
