package middleware

import (
	"net/http"

	"time"

	"github.com/vsco/dcdr/client"
)

const (
	// IfNoneMatchHeader http header containing the clients CurrentSHA
	IfNoneMatchHeader = "If-None-Match"
	// EtagHeader header used to pass the CurrentSHA in responses
	EtagHeader = "Etag"
	// LastModified date when the feature set was last updated
	LastModifiedHeader = "Last-Modified"
	// CacheControlHeader sets the caches control header for the response
	CacheControlHeader = "Cache-Control"
	// CacheControl ensure no client-side or proxy caching
	CacheControl = "private, max-age=0, no-cache, no-store, must-revalidate, proxy-revalidate"
	// PragmaHeader sets the pragma header for the response
	PragmaHeader = "Pragma"
	// Pragma ensure no client-side caching
	Pragma = "no-cache"
	// ExpiresHeader sets the expires header for the response
	ExpiresHeader = "Expires"
	// Expires ensure no client-side caching
	Expires = "0"
)

// NotModified checks the requests If-None-Match header for a matching
// CurrentSHA in the Client.
func NotModified(sha string, r *http.Request) bool {
	if val := r.Header.Get(IfNoneMatchHeader); val != "" && sha != "" {
		return sha == val
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

// HTTPCachingHandler middleware that provides HTTP level caching
// using the If-None-Match header. If the value of this header contains
// a matching CurrentSHA this handler will write a 304 status and return.
func HTTPCachingHandler(dcdr client.IFace) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if sha := dcdr.Info().CurrentSHA; sha != "" {
				w.Header().Set(EtagHeader, sha)
			}

			if date := dcdr.Info().LastModifiedDate; date != 0 {
				lmd := time.Unix(date, 0).Format(time.RFC1123)
				w.Header().Set(LastModifiedHeader, lmd)
			}

			SetCacheHeaders(w, r)

			if NotModified(dcdr.Info().CurrentSHA, r) {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
