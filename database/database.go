package database

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	InitQuery = `CREATE TABLE IF NOT EXISTS messages (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        room varchar(255) DEFAULT NULL,
        author varchar(255) DEFAULT NULL,
        message text,
        at datetime DEFAULT NULL,
        created_at datetime NOT NULL,
        updated_at datetime NOT NULL
    );`

	InsertionQuery                = "INSERT INTO messages (room, author, message, at, created_at, updated_at) VALUES (:room, :author, :message, :at, datetime('now'), datetime('now'))"
	SelectAllRoomsQuery           = "SELECT DISTINCT room FROM messages ORDER BY room"
	SelectLatestQuery             = "SELECT * FROM messages ORDER BY at DESC LIMIT 0,?"
	SelectLatestByRoomQuery       = "SELECT * FROM messages WHERE room = ? ORDER BY at DESC LIMIT 0,?"
	SelectLatestByAuthorQuery     = "SELECT * FROM messages WHERE author = ? ORDER BY at DESC LIMIT 0,?"
	SelectPriorMessagesQuery      = "SELECT * FROM messages WHERE room = ? AND at < ? ORDER BY at DESC LIMIT 0,?"
	SelectSubsequentMessagesQuery = "SELECT * FROM messages WHERE room = ? AND at > ? ORDER BY at ASC LIMIT 0,?"
	SelectMessageByIdQuery        = "SELECT * FROM messages WHERE id = ?"
	SelectByFuzzyMessageQuery     = "SELECT * FROM messages WHERE message LIKE ? ORDER BY id DESC"
)

type DB struct {
	Connection *sqlx.DB
}

var ErrInvalidQuery = fmt.Errorf("query is invalid")

func databaseURL() string {
	return os.Getenv("GOSSIP_DB_PATH")
}

func New() *DB {
	db := &DB{}
	_ = db.InitDB(context.Background())
	return db
}

func (db *DB) InitDB(ctx context.Context) error {
	_, err := db.GetConnection().ExecContext(ctx, InitQuery)
	return err
}

func (db *DB) Connect(ctx context.Context) (*sqlx.DB, error) {
	if db.Connection == nil {
		conn, err := sqlx.ConnectContext(ctx, "sqlite3", databaseURL())
		if err != nil {
			return nil, err
		}
		db.Connection = conn
	}

	return db.Connection, nil
}

func (db *DB) GetConnection() *sqlx.DB {
	if db.Connection == nil {
		db.Connection = sqlx.MustConnect("sqlite3", databaseURL())
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

func (db *DB) AllRooms() ([]string, error) {
	allRooms := []string{}
	err := db.GetConnection().Select(&allRooms, SelectAllRoomsQuery)
	return allRooms, err
}

func (db *DB) Find(id int) (*Message, error) {
	msg := &Message{}
	err := db.GetConnection().Get(msg, SelectMessageByIdQuery, id)
	return msg, err
}

func (db *DB) PriorMessages(room, at string, limit int) ([]Message, error) {
	messages := SortableMessages{}
	err := db.GetConnection().Select(&messages, SelectPriorMessagesQuery, room, at, limit)
	sort.Stable(messages)
	return []Message(messages), err
}

func (db *DB) SubsequentMessages(room, at string, limit int) ([]Message, error) {
	messages := []Message{}
	err := db.GetConnection().Select(&messages, SelectSubsequentMessagesQuery, room, at, limit)
	return messages, err
}

func (db *DB) LatestMessages(limit int) ([]Message, error) {
	messages := []Message{}
	err := db.GetConnection().Select(&messages, SelectLatestQuery, limit)
	return messages, err
}

func (db *DB) LatestMessagesByRoom(room string, limit int) ([]Message, error) {
	messages := []Message{}
	err := db.GetConnection().Select(&messages, SelectLatestByRoomQuery, room, limit)
	return messages, err
}

func (db *DB) LatestMessagesByAuthor(author string, limit int) ([]Message, error) {
	messages := []Message{}
	err := db.GetConnection().Select(&messages, SelectLatestByAuthorQuery, author, limit)
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

func (db *DB) ListByFuzzyMessage(searchTerm string) ([]Message, error) {
	messages := []Message{}
	if searchTerm == "" {
		return messages, ErrInvalidQuery
	}
	err := db.GetConnection().Select(&messages, SelectByFuzzyMessageQuery, "%"+searchTerm+"%")
	return messages, err
}
