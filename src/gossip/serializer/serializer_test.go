package serializer

import (
	"reflect"
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

func TestMarshalJsonError(t *testing.T) {
	something := func() string {
		return "Error, I say!"
	}
	expected := "json: unsupported type: func() string"
	actual := MarshalJson(something)
	if actual != expected {
		t.Fatalf("Marshalling failed: got '%s' instead of '%s'", actual, expected)
	}
}

func TestParseJavaScriptTime(t *testing.T) {
	thyme := "Sun, 04 Jan 2015 04:26:57 GMT"
	loc, _ := time.LoadLocation("UTC")
	expected := time.Date(2015, time.January, 4, 4, 26, 57, 0, loc)
	actual := ParseJavaScriptTime(thyme).UTC()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("ParseJavaScriptTime() failed: expected '%s', got '%s'", expected, actual)
	}
}

func TestParseJavaScriptTimeError(t *testing.T) {
	thyme := "Sat, 50 Jan 2015 15:10:12 PST"
	actual := ParseJavaScriptTime(thyme)
	if actual != nil {
		t.Fatalf("ParseJavaScriptTime() failed: expected nil, got '%s'", actual)
	}
}

func TestTimeToXML(t *testing.T) {
	actual := TimeToXML(time.Unix(0, 0))
	expected := "1970-01-01 00:00:00 +0000"
	if actual != expected {
		t.Fatalf("TimeToXML() failed: expected '%s', got '%s'", expected, actual)
	}
}
