package main

import (
	"encoding/json"
	"fmt"
)

type ResponseMessage struct {
	OK     string    `json:"ok"`
	Limit  string    `json:"limit"`
	Values []Message `json:"values"`
	Error  error     `json:"error"`
}

func (r *ResponseMessage) Json() string {
	resp_json, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error converting ResponseMessage to JSON", r, err)
		return ""
	}
	return string(resp_json)
}

func basicResponseMessage(okVal string) string {
	r := &ResponseMessage{}
	r.OK = okVal
	return r.Json()
}

func messagesResponseMessage(limit string, messages []Message) string {
	r := &ResponseMessage{}
	r.OK = "true"
	r.Limit = limit
	r.Values = messages
	return r.Json()
}

func errorMessage(err error) string {
	r := &ResponseMessage{}
	r.OK = "false"
	r.Error = err
	return r.Json()
}
