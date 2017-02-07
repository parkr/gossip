package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sniffer := newResponseWriterSniffer(w)
		logForReq(r, fmt.Sprintf("Started %s %q from %s", r.Method, r.URL.String(), r.RemoteAddr))
		start := time.Now()
		h.ServeHTTP(sniffer, r)
		logForReq(r, fmt.Sprintf("Returning %d in %s", sniffer.Code(), time.Since(start)))
	})
}

func recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("%q", err), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var bind string
	flag.StringVar(&bind, "bind", ":"+port, "Thing the server should bind to when it runs. Usually a port.")
	flag.Parse()

	http.Handle("/", recoverMiddleware(requestIDMiddleware(loggingMiddleware(handler))))

	log.Println("Launching gossip on", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatalln("server crashed:", err)
	}
}
