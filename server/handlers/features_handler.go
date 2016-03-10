package handlers

import (
	"net/http"

	"strings"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
)

const DcdrScopesHeader = "x-dcdr-scopes"
const CurrentSha = "current_sha"
const ContentTypeHeader = "Content-Type"
const ContentType = "application/json"

func GetScopes(r *http.Request) []string {
	return strings.Split(r.Header.Get(DcdrScopesHeader), ",")
}

func SetResponseHeaders(w http.ResponseWriter) {
	w.Header().Set(ContentTypeHeader, ContentType)
}

func FeaturesHandler(cfg *config.Config, c client.ClientIFace) func(
	w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		json, err := c.WithScopes(GetScopes(r)...).ScopedMap().ToJson()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetResponseHeaders(w)
		w.Write(json)
	}
}
