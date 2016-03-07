package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/vsco/dcdr/client"
)

var FixturePath, _ = filepath.Abs("../../config/decider_fixtures.json")

func renderFeatures(c *client.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		scope := fmt.Sprintf("cc/%s", strings.ToLower(r.Header.Get("X-Country")))
		scoped := c.WithScopes(scope)

		js, err := json.MarshalIndent(scoped.Features(), "", "  ")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func main() {
	cfg := &client.Config{
		WatchPath: FixturePath,
	}

	c, err := client.New(cfg).Watch()

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", renderFeatures(c))
	http.ListenAndServe(":3000", nil)
}
