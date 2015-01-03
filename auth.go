package main

import (
	"github.com/martini-contrib/auth"
	"log"
	"net/http"
	"os"
)

func TokenAuthHandler(h http.Handler) http.Handler {
	token := os.Getenv("GOSSIP_AUTH_TOKEN")
	handler := func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.Fatal("Error processing the form!!!")
			log.Fatal(err)
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
