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
	Code   int       `json:"code"`
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
	if okVal == "true" {
		r.Code = 200
	} else {
		r.Code = 500
	}
	return r.Json()
}

func messagesResponseMessage(limit string, messages []Message) string {
	r := &ResponseMessage{}
	r.OK = "true"
	r.Code = 200
	r.Limit = limit
	r.Values = messages
	return r.Json()
}

func errorMessage(err error) (int, string) {
	r := &ResponseMessage{}
	r.OK = "false"
	r.Code = 500
	r.Error = err
	return r.Code, r.Json()
}
