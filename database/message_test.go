package database

import "testing"

func TestString(t *testing.T) {
	msg := &Message{
		Room:    "#gossip",
		Author:  "parkr",
		Message: "You found your way, skipper.",
		At:      "Mon, 02 Jan 2006 15:04:05 MST",
	}
	expected := "<#gossip by parkr at Mon, 02 Jan 2006 15:04:05 MST: You found your way, skipper.>"
	actual := msg.String()
	if actual != expected {
		t.Fatalf("String() failed: expected '%s', got '%s'", expected, actual)
	}
}

func Test_CreatedAtRFC3339(t *testing.T) {
	db := New()
	defer db.Close()
	prepareDatabase(t, db)

	msgs, err := db.LatestMessages(1)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	createdAt := msgs[0].CreatedAtRFC3339()
	if createdAt == "" {
		t.Fatal("expected created at to return a time, but didn't")
	}
}
