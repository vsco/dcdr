package handlers

import (
	"net/http"

	"strings"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/models"
)

const (
	// DcdrScopesHeader comma delimited scopes to pass to the client
	DcdrScopesHeader = "x-dcdr-scopes"
	// ContentTypeHeader header for content type
	ContentTypeHeader = "Content-Type"
	// ContentType set JSON content type for responses
	ContentType = "application/json"
)

// SetResponseHeaders set the common response headers
func SetResponseHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentTypeHeader, ContentType)
	w.Header().Set(DcdrScopesHeader, r.Header.Get(DcdrScopesHeader))
}

// GetScopes parses the comma delimited string from DcdrScopesHeader into
// a slice of strings.
//
// x-dcdr-scopes: "a/b/c, d" => []string{"a/b/c", "d"}
func GetScopes(r *http.Request) []string {
	scopes := strings.Split(r.Header.Get(DcdrScopesHeader), ",")
	for i := 0; i < len(scopes); i++ {
		scopes[i] = strings.TrimSpace(scopes[i])
	}

	return scopes
}

// ScopeMapFromRequest helper method for returning a FeatureMap scoped to
// the values found in DcdrScopesHeader.
func ScopeMapFromRequest(c client.IFace, r *http.Request) *models.FeatureMap {
	return c.WithScopes(GetScopes(r)...).ScopedMap()
}

// FeaturesHandler default handler for serving a FeatureMap via HTTP
func FeaturesHandler(c client.IFace) func(
	w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json, err := ScopeMapFromRequest(c, r).ToJSON()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetResponseHeaders(w, r)
		w.Write(json)
	}
}
