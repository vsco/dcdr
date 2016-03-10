package main

import (
	"net/http"

	"github.com/vsco/dcdr/server"
	"github.com/zenazn/goji/web"
)

const AuthorizationHeader = "Authorization"

// MockAuth example authentication middleware.
// Checks for any value in the http Authorization header.
// If no value is found a 401 status is sent.
func MockAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(AuthorizationHeader) != "" {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}

	return http.HandlerFunc(fn)
}

func main() {
	// Create a new Server and Client
	srv := server.NewDefault()

	// Add the MockAuth to the middleware chain
	srv.Use(MockAuth)

	// Begin serving on :8000
	// curl -sH "Authorization: authorized" :8000/dcdr.json
	srv.Serve()
}
