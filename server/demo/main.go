package main

import (
	"net/http"

	"fmt"
	"strings"

	"log"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/server"
	"github.com/vsco/dcdr/server/handlers"
)

const AuthorizationHeader = "Authorization"
const CountryCodeHeader = "X-Country"

// MockAuth example authentication middleware.
// Checks for any value in the http Authorization header.
// If no value is found a 401 status is sent.
func MockAuth(c client.IFace) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(AuthorizationHeader) != "" {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}

		return http.HandlerFunc(fn)
	}
}

// ScopedCountryCode appends the country-codes/xx scope as the X-Country header
func ScopedCountryCode(c client.IFace) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cc := strings.ToLower(r.Header.Get(CountryCodeHeader))

			if cc != "" {
				handlers.AppendScope(r, fmt.Sprintf("country-codes/%s", cc))
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func main() {
	// Create a new Server and Client
	srv, err := server.NewDefault()

	if err != nil {
		log.Fatal(err)
	}

	// Add the MockAuth and ScopedCountryCode to the middleware chain
	srv.Use(MockAuth, ScopedCountryCode)

	// Begin serving on :8000
	// curl -sH "Authorization: authorized" :8000/dcdr.json
	srv.Serve()
}
