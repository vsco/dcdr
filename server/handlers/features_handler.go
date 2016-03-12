package handlers

import (
	"net/http"

	"strings"

	"github.com/vsco/dcdr/client"
)

const DcdrScopesHeader = "x-dcdr-scopes"
const IfNoneMatchHeader = "If-None-Match"
const EtagHeader = "Etag"
const ContentTypeHeader = "Content-Type"
const CacheControlHeader = "Cache-Control"
const PragmaHeader = "Pragma"
const ExpiresHeader = "Expires"

const ContentType = "application/json"
const CacheControl = "no-cache, no-store, must-revalidate"
const Pragma = "no-cache"
const Expires = "0"

func GetScopes(r *http.Request) []string {
	return strings.Split(r.Header.Get(DcdrScopesHeader), ",")
}

func SetResponseHeaders(w http.ResponseWriter, r *http.Request, sha string) {
	w.Header().Set(ContentTypeHeader, ContentType)
	w.Header().Set(EtagHeader, sha)
	w.Header().Set(CacheControlHeader, CacheControl)
	w.Header().Set(PragmaHeader, Pragma)
	w.Header().Set(ExpiresHeader, Expires)
	w.Header().Set(DcdrScopesHeader, r.Header.Get(DcdrScopesHeader))
}

func NotModified(sha string, r *http.Request) bool {
	if v := r.Header.Get(IfNoneMatchHeader); v != "" && sha != "" {
		return sha == r.Header.Get(IfNoneMatchHeader)
	}

	return false
}

func FeaturesHandler(c client.ClientIFace) func(
	w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if NotModified(c.CurrentSha(), r) {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		json, err := c.WithScopes(GetScopes(r)...).ScopedMap().ToJson()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		SetResponseHeaders(w, r, c.CurrentSha())
		w.Write(json)
	}
}
