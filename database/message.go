package database

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int    `json:"id" db:"id"`
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

func (msg *Message) CreatedAtRFC3339() string {
	return sqlTimeToGoTime(msg.CreatedAt).Format(time.RFC3339)
}

func (msg *Message) GetUpdatedAt() time.Time {
	return sqlTimeToGoTime(msg.UpdatedAt)
}

func sqlTimeToGoTime(sqlTime string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", sqlTime)
	if err != nil {
		panic(err)
	}
	return t
}

type SortableMessages []Message

func (m SortableMessages) Len() int {
	return len(m)
}

func (m SortableMessages) Less(i, j int) bool {
	return m[i].At < m[j].At
}

func (m SortableMessages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
