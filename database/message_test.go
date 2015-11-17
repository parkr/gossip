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
