package gossip

import (
	"net/http"
)

type sniffer struct {
	StatusCode int

	w http.ResponseWriter
}

func NewResponseWriterSniffer(w http.ResponseWriter) *sniffer {
	return &sniffer{
		StatusCode: -1,
		w:          w,
	}
}

func (s *sniffer) Header() http.Header {
	return s.w.Header()
}

func (s *sniffer) Write(data []byte) (int, error) {
	return s.w.Write(data)
}

func (s *sniffer) WriteHeader(code int) {
	s.w.WriteHeader(code)
	s.StatusCode = code
}

// Status code to return.
func (s *sniffer) Code() int {
	if s.StatusCode == -1 {
		return 200
	}
	return s.StatusCode
}
