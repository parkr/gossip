package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martini-contrib/auth"
)

var authToken = os.Getenv("GOSSIP_AUTH_TOKEN")

func getAuthToken() string {
	if authToken != "" {
		return authToken
	}
	panic("GOSSIP_AUTH_TOKEN is a required environment variable")
}

func TokenAuthHandler(h http.Handler) http.Handler {
	token := getAuthToken()
	handler := func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.Println("Error processing the form:", err)
		}
		givenToken := req.Form.Get("access_token")
		if req.URL.Path != "/" && !auth.SecureCompare(givenToken, token) {
			res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		} else {
			h.ServeHTTP(res, req)
		}
	}
	return http.HandlerFunc(handler)
}
