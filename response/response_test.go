package response

import (
	"github.com/parkr/gossip/database"

	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	newResponse := New()
	if newResponse == nil {
		t.Fatalf("New() didn't create a new ResponseMessage")
	}
}

func TestWithError(t *testing.T) {
	res := New().WithError(errors.New("Broken."))
	if res.Error == nil || res.Error.Error() != "Broken." {
		t.Fatalf("WithError() failed: expected the error to be set, got '%s'", res.Error)
	}
}

func TestWithLimit(t *testing.T) {
	res := New().WithLimit("10")
	if res.Limit != "10" {
		t.Fatalf("WithLimit() failed: got '%s' instead of 10", res.Limit)
	}
}

func TestWithLimitOverriding(t *testing.T) {
	res := New().WithLimit("10").WithLimit("100")
	if res.Limit != "100" {
		t.Fatalf("WithLimit() failed: got '%s' instead of 100", res.Limit)
	}
}

func TestWithMessages(t *testing.T) {
	res := New()
	if res.Messages != nil {
		t.Fatal("WithMessages() can't rely on a nil setting of `Messages`")
	}
	res.WithMessages([]database.Message{})
	if res.Messages == nil {
		t.Fatalf("WithMessages() failed: got '%s' instead of a new slice", res.Messages)
	}
}

func TestWithMessageWithNoMessages(t *testing.T) {
	res := New()
	msg := database.Message{}
	res.WithMessage(&msg)
	if !reflect.DeepEqual(res.Messages[0], msg) {
		t.Fatalf("WithMessage() failed: got '%s' instead of %s", &res.Messages[0], msg)
	}
}

func TestWithMessageWithOtherMessages(t *testing.T) {
	res := New().WithMessages([]database.Message{database.Message{Room: "#gossip"}})
	msg := database.Message{Room: "#github"}
	res.WithMessage(&msg)
	if !reflect.DeepEqual(res.Messages[1], msg) {
		t.Fatalf("WithMessage() failed: got '%s' instead of %s", res.Messages[1], msg)
	}
}

func TestChaining(t *testing.T) {
	res := New()
	if res.Messages != nil {
		t.Fatalf("Default value of .Messages is not nil, but rather '%s'", res.Messages)
	}
	if res.Limit != "" {
		t.Fatalf("Default value of .Limit is not '', but rather '%s'", res.Limit)
	}
	res.WithLimit("5").WithMessages([]database.Message{})
	if res.Limit != "5" {
		t.Fatalf("WithLimit() failed: got '%s' instead of 5", res.Limit)
	}
	if res.Messages == nil {
		t.Fatalf("WithMessages() failed: got '%s' instead of a new slice", res.Messages)
	}
}

func TestJson(t *testing.T) {
	res := New().WithLimit("5").WithMessages([]database.Message{})
	json := res.Json()
	expected := `{"messages":[],"limit":"5","error":null}`
	if json != expected {
		t.Fatalf("Json() failed: got '%s' instead of '%s'", json, expected)
	}
}
