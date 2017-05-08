package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

var reqIDKey = "reqID"
var prefix string
var reqid uint64

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func newRequestID() string {
	myid := atomic.AddUint64(&reqid, 1)
	return fmt.Sprintf("%s-%06d", prefix, myid)
}

func contextWithRequestID(r *http.Request) context.Context {
	return context.WithValue(r.Context(), reqIDKey, newRequestID())
}

func logForReq(r *http.Request, message string) {
	log.Printf("[%s] %s", requestID(r), message)
}

// Fetches the request ID from the request's Context.
func requestID(r *http.Request) string {
	reqID := r.Context().Value(reqIDKey)
	switch v := reqID.(type) {
	case string:
		return v
	default:
		log.Println("Haven't the slightest idea what request ID this is:", v)
		return "NOT-FOUND"
	}
}

func requestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq := r.WithContext(contextWithRequestID(r))
		w.Header().Set("X-Request-Id", requestID(newReq))
		h.ServeHTTP(w, newReq)
	})
}
