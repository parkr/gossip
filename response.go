package main

import (
	"encoding/json"
	"fmt"
)

type ResponseMessage struct {
	OK     string    `json:"ok"`
	Code   int       `json:"code"`
	Values []Message `json:"values"`
	Limit  string    `json:"limit"`
	Error  error     `json:"error"`
}

func (r *ResponseMessage) Json() string {
	resp_json, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error converting ResponseMessage to JSON", r, err)
		return err.Error()
	}
	return string(resp_json)
}

func (r *ResponseMessage) MartiniResp() (int, string) {
	return r.Code, r.Json()
}

func singleMessageResponse(msg Message) (int, string) {
	msgs := []Message{msg}
	return messagesResponse("", msgs)
}

func messagesResponse(limit string, messages []Message) (int, string) {
	r := &ResponseMessage{
		"true",
		200,
		messages,
		limit,
		nil,
	}
	return r.MartiniResp()
}

func errorResponse(code int, err error) (int, string) {
	r := &ResponseMessage{
		"false",
		code,
		nil,
		"",
		err,
	}
	return r.MartiniResp()
}

func internalErrorResponse(err error) (int, string) {
	return errorResponse(500, err)
}
