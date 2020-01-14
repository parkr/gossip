package database

import (
	"github.com/parkr/gossip/serializer"

	"os"
	"testing"
)

func withDatabaseAuthenticationInfo(user, password, hostname, dbname string, f func()) {
	cachedDatabaseURL = ""
	oldGossipDbUsername := os.Getenv("GOSSIP_DB_USERNAME")
	os.Setenv("GOSSIP_DB_USERNAME", user)
	oldGossipDbPassword := os.Getenv("GOSSIP_DB_PASSWORD")
	os.Setenv("GOSSIP_DB_PASSWORD", password)
	oldGossipDbHostname := os.Getenv("GOSSIP_DB_HOSTNAME")
	os.Setenv("GOSSIP_DB_HOSTNAME", "")
	oldGossipDbName := os.Getenv("GOSSIP_DB_DBNAME")
	os.Setenv("GOSSIP_DB_DBNAME", dbname)
	f()
	cachedDatabaseURL = ""
	os.Setenv("GOSSIP_DB_USERNAME", oldGossipDbUsername)
	os.Setenv("GOSSIP_DB_PASSWORD", oldGossipDbPassword)
	os.Setenv("GOSSIP_DB_HOSTNAME", oldGossipDbHostname)
	os.Setenv("GOSSIP_DB_DBNAME", oldGossipDbName)
}

func TestDatabaseURL(t *testing.T) {
	withDatabaseAuthenticationInfo("travis", "blah", "", "gossip_test", func() {
		actual := databaseURL()
		expected := "travis:blah@/gossip_test"
		if actual != expected {
			t.Fatalf("databaseURL() failed: expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestNew(t *testing.T) {
	db := New()

	if db == nil {
		t.Fatal("New() failed: expected a db, got nil")
	}

	if db.Connection != nil {
		t.Fatalf("New() failed: expected the connection to be nil, got %+v", db.Connection)
	}
}

func TestGetConnection(t *testing.T) {
	db := New()
	if db.GetConnection() == nil {
		t.Fatal("GetConnection() failed: expected the connection to exist, got nil")
	}
}

func TestClose(t *testing.T) {
	db := New()
	err := db.Close()
	if err != nil {
		t.Fatalf("Close() failed: encountered error '%s'", err)
	}

	if db.Connection != nil {
		t.Fatalf("Close() failed: .Connection should be nil, but is '%+v'", db.Connection)
	}
}

func TestLatestMessages(t *testing.T) {
	db := New()
	defer db.Close()

	msgs, err := db.LatestMessages(1)

	if err != nil {
		t.Fatalf("LatestMessages() failed: encountered error '%s'", err)
	}

	if &msgs == nil {
		t.Fatal("LatestMessages() failed: expected a []Message, got nil")
	}
}

func TestInsertMessage(t *testing.T) {
	db := New()
	defer db.Close()

	msg := map[string]interface{}{
		"room":    "#jekyll",
		"author":  "parker",
		"message": "the meeeessaaaaggeeee",
		"at":      serializer.ParseJavaScriptTime("Mon, 02 Jan 2006 15:04:05 MST"),
	}

	actual, err := db.InsertMessage(msg)

	if err != nil {
		t.Fatalf("InsertMessage() failed: encountered error '%s'", err)
	}

	if &actual == nil {
		t.Fatal("InsertMessage() failed: expected a Message, got nil")
	}

	if actual.Room != msg["room"] {
		t.Fatalf("InsertMessage() failed: expected .Room to be '%s', got '%s'", msg["room"], actual.Room)
	}

	if actual.Author != msg["author"] {
		t.Fatalf("InsertMessage() failed: expected .Author to be '%s', got '%s'", msg["author"], actual.Author)
	}

	if actual.Message != msg["message"] {
		t.Fatalf("InsertMessage() failed: expected .Message to be '%s', got '%s'", msg["message"], actual.Message)
	}

	expectedAt := "2006-01-02 15:04:05"
	if actual.At != expectedAt {
		t.Fatalf("InsertMessage() failed: expected .At to be '%s', got '%s'", expectedAt, actual.At)
	}
}

func TestInsertMessageError(t *testing.T) {
	db := New()
	defer db.Close()

	msg := map[string]interface{}{
		"author":  "parker",
		"message": "the meeeessaaaaggeeee",
		"at":      serializer.ParseJavaScriptTime("Mon, 02 Jan 2006 15:04:05 UTC"),
	}

	actual, err := db.InsertMessage(msg)

	if err == nil {
		t.Fatal("InsertMessage() failed: expected error but got nil")
	}

	if actual != nil {
		t.Fatalf("InsertMessage() failed: expected nil, got %s", actual)
	}
}

func TestFind(t *testing.T) {
	db := New()
	defer db.Close()

	msg, err := db.Find(1)
	if err != nil {
		t.Fatalf("Find() failed: encountered error '%s'", err)
	}

	if &msg == nil {
		t.Fatal("Find() failed: expected a message, got nil")
	}
}
