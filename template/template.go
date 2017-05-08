// template defines gossip templates and their data structures.
package template

import (
	"html/template"
	"strings"
)

var addtlFuncs = template.FuncMap{
	"encodeRoom": encodeRoom,
}

func encodeRoom(room string) string {
	return strings.Replace(room, "#", "%23", -1)
}
