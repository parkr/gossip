package main

import (
	"github.com/martini-contrib/auth"
	"log"
	"net/http"
	"os"
)

func TokenAuthHandler() http.HandlerFunc {
	token := os.Getenv("GOSSIP_AUTH_TOKEN")
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.Fatal("Error processing the form!!!")
			log.Fatal(err)
		}
		givenToken := req.Form.Get("access_token")
		if !auth.SecureCompare(givenToken, token) {
			res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}
