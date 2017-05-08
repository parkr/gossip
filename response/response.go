package response

import (
	"github.com/parkr/gossip/database"
	"github.com/parkr/gossip/serializer"
)

type ResponseMessage struct {
	Messages []database.Message `json:"messages"`
	Limit    int                `json:"limit"`
	Error    error              `json:"error"`
}

func New() *ResponseMessage {
	return &ResponseMessage{}
}

func (r *ResponseMessage) WithError(err error) *ResponseMessage {
	r.Error = err
	return r
}

func (r *ResponseMessage) WithLimit(limit int) *ResponseMessage {
	r.Limit = limit
	return r
}

func (r *ResponseMessage) WithMessages(messages []database.Message) *ResponseMessage {
	r.Messages = messages
	return r
}

func (r *ResponseMessage) WithMessage(message *database.Message) *ResponseMessage {
	if r.Messages == nil {
		r.Messages = []database.Message{*message}
	} else {
		r.Messages = append(r.Messages, *message)
	}
	return r
}

func (r *ResponseMessage) Json() string {
	return serializer.MarshalJson(r)
}
