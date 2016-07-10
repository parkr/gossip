package main

import (
	"log"
	"net/http"
	"regexp"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func serve() {
	goji.DefaultMux.Compile()
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", goji.DefaultMux)

	listener := bind.Socket(bind.Sniff())
	log.Println("Starting Goji on", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() {
		log.Printf("Goji stopped")
		log.Printf("Shutting down the server")
		handler.DB.Close()
		log.Printf("Database shut down. Terminating the process.")
	})

	err := graceful.Serve(listener, http.DefaultServeMux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}

func main() {
	goji.Get("/", handler.SayHello)
	byRoomPattern := regexp.MustCompile(`^/room/(?P<room>[%a-zA-Z0-9]+)$`)
	goji.Get(byRoomPattern, handler.FindMessagesByRoom)
	byAuthorPattern := regexp.MustCompile(`^/messages/by/(?P<author>[a-zA-Z0-9_-]+)$`)
	goji.Get(byAuthorPattern, handler.FindMessagesByAuthor)
	messageContextPatther := regexp.MustCompile(`^/messages/(?P<id>[0-9]+)/context$`)
	goji.Get(messageContextPatther, handler.FindMessageWithContext)

	messages := web.New()
	messages.Use(TokenAuthHandler)

	byIDPattern := regexp.MustCompile(`^(?P<id>[0-9]+)$`)
	messages.Get(byIDPattern, handler.FindMessageById)
	messages.Get("/latest", handler.FetchLatestMessages)
	messages.Post("/log", handler.StoreMessage)

	goji.Handle("/api/messages/*", messages)

	serve()
}
