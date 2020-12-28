package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/parkr/gossip"
	_ "github.com/parkr/gossip/statik"
	"github.com/rakyll/statik/fs"
)

func remoteAddr(r *http.Request) string {
	if val := r.Header.Get("X-Forwarded-For"); val != "" {
		return val
	}
	return r.RemoteAddr
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sniffer := gossip.NewResponseWriterSniffer(w)
		gossip.LogWithRequestID(r, fmt.Sprintf("Started %s %q from %s", r.Method, r.URL.String(), remoteAddr(r)))
		start := time.Now()
		h.ServeHTTP(sniffer, r)
		gossip.LogWithRequestID(r, fmt.Sprintf("Returning %d in %s", sniffer.Code(), time.Since(start)))
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

	handler := gossip.NewHandler()

	statikFS, err := fs.New()
	if err != nil {
		log.Printf("error creating statik FS: %+v", err)
	} else {
		http.Handle("/assets/", recoverMiddleware(gossip.RequestIDMiddleware(loggingMiddleware(http.FileServer(statikFS)))))
	}
	http.Handle("/", recoverMiddleware(gossip.RequestIDMiddleware(loggingMiddleware(handler))))

	log.Println("Launching gossip on", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatalln("server crashed:", err)
	}
	handler.DB.Close()
}
