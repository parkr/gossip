package main

import (
	"errors"
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

func authenticate(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	givenToken := r.Form.Get("access_token")
	if r.URL.Path != "/" && !auth.SecureCompare(givenToken, getAuthToken()) {
		return errors.New("no access!!!")
	}

	return nil
}

func TokenAuthHandler(h http.Handler) http.Handler {
	handler := func(res http.ResponseWriter, req *http.Request) {
		if err := authenticate(req); err != nil {
			http.Error(res, "401 Not Authorized", http.StatusUnauthorized)
		} else {
			h.ServeHTTP(res, req)
		}
	}
	return http.HandlerFunc(handler)
}
