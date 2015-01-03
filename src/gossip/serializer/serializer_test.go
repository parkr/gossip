package serializer

import (
	"testing"
	"time"
)

func TestMarshalJson(t *testing.T) {
	something := map[string]interface{}{
		"host":  123,
		"port":  "locomotion",
		"death": time.Unix(10000000, 0),
	}
	expected := "{\"death\":\"1970-04-26T17:46:40Z\",\"host\":123,\"port\":\"locomotion\"}"
	actual := MarshalJson(something)
	if actual != expected {
		t.Fatalf("Marshalling failed: got '%s' instead of '%s'", actual, expected)
	}
}
