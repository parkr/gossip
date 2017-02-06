package main

import (
	"flag"
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

	listener := bind.Default()
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
	flag.Parse()

	goji.Get("/", handler.SayHello)

	messages := web.New()
	messages.Use(TokenAuthHandler)

	pattern := regexp.MustCompile(`^(?P<id>[0-9]+)$`)
	messages.Get(pattern, handler.FindMessageById)
	messages.Get("/latest", handler.FetchLatestMessages)
	messages.Post("/log", handler.StoreMessage)

	goji.Handle("/api/messages/*", messages)

	serve()
}
