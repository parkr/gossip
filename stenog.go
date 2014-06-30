package main

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
)

func main() {
	// Setup Martini
	m := martini.Classic()
	h := &Handler{newDB()}
	m.Get("/", h.SayHello)
	m.Group("/api/messages", func(r martini.Router) {
		r.Get("/(?P<id>[0-9]+)", TokenAuthHandler(), h.FindMessageById)
		r.Get("/latest", TokenAuthHandler(), h.FetchLatestMessages)
		r.Post("/log", TokenAuthHandler(), binding.Bind(Message{}), h.StoreMessage)
	})
	m.Run()
}
