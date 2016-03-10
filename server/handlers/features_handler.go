package handlers

import (
	"net/http"

	"strings"

	"github.com/vsco/dcdr/client"
)

const DcdrScopesHeader = "x-dcdr-scopes"
const ContentTypeHeader = "Content-Type"
const ContentType = "application/json"

func GetScopes(r *http.Request) []string {
	return strings.Split(r.Header.Get(DcdrScopesHeader), ",")
}

func SetResponseHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentTypeHeader, ContentType)
	w.Header().Set(DcdrScopesHeader, r.Header.Get(DcdrScopesHeader))
}

func FeaturesHandler(c client.ClientIFace) func(
	w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		json, err := c.WithScopes(GetScopes(r)...).ScopedMap().ToJson()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetResponseHeaders(w, r)
		w.Write(json)
	}
}
