package middleware

import (
	"net/http"

	"github.com/vsco/dcdr/client"
	"github.com/zenazn/goji/web"
)

const (
	// IfNoneMatchHeader http header containing the clients CurrentSha
	IfNoneMatchHeader = "If-None-Match"
	// EtagHeader header used to pass the CurrentSha in responses
	EtagHeader = "Etag"
	// CacheControlHeader sets the caches control header for the response
	CacheControlHeader = "Cache-Control"
	// PragmaHeader sets the pragma header for the response
	PragmaHeader = "Pragma"
	// ExpiresHeader sets the expires header for the response
	ExpiresHeader = "Expires"
	// CacheControl ensure no client-side caching
	CacheControl = "no-cache, no-store, must-revalidate"
	// Pragma ensure no client-side caching
	Pragma = "no-cache"
	// Expires ensure no client-side caching
	Expires = "0"
)

// NotModified checks the requests If-None-Match header for a matching
// CurrentSha in the Client.
func NotModified(sha string, r *http.Request) bool {
	if v := r.Header.Get(IfNoneMatchHeader); v != "" && sha != "" {
		return sha == r.Header.Get(IfNoneMatchHeader)
	}

	return false
}

// SetCacheHeaders ensure response is not cached on the client.
// Caching should be done using the Etag and If-None-Match headers.
func SetCacheHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CacheControlHeader, CacheControl)
	w.Header().Set(PragmaHeader, Pragma)
	w.Header().Set(ExpiresHeader, Expires)
}

// HTTPCachingHandle middleware that provides HTTP level caching
// using the If-None-Match header. If the value of this header contains
// a matching CurrentSha this handler will write a 304 status and return.
func HTTPCachingHandler(dcdr client.IFace) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if sha := dcdr.CurrentSha(); sha != "" {
				w.Header().Set(EtagHeader, sha)
			}

			SetCacheHeaders(w, r)

			if NotModified(dcdr.CurrentSha(), r) {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
