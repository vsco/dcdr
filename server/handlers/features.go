package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/vsco/dcdr/watcher"
	"github.com/zenazn/goji/web"
)

const (
	IfNoneMatchHeader  = "If-None-Match"
	CacheHeaders       = "no-cache, no-store, must-revalidate"
	CacheControlHeader = "Cache-Control"
	EtagHeader         = "Etag"
	ContentTypeHeader  = "Content-Type"
	ContentType        = "application/json"
)

type FeatureHandler struct {
	Watcher *watcher.Watcher
	cfg     *watcher.Config
}

func NewFeatureHandler(l *watcher.Watcher, cfg *watcher.Config) (fh *FeatureHandler) {
	fh = &FeatureHandler{
		Watcher: l,
		cfg:     cfg,
	}

	return
}

func (fh *FeatureHandler) Serve(c web.C, w http.ResponseWriter, r *http.Request) {
	fh.AddHeaders(w)

	if etag, ok := r.Header[IfNoneMatchHeader]; ok {
		if !fh.Watcher.Expired(etag[0]) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	ctxFts := watcher.GetFeatureMap(&c)
	features, err := fh.Watcher.MergeFeatures(&ctxFts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	bts, err := json.Marshal(features)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(bts)
}

// AddHeaders sets cache control headers
func (fh *FeatureHandler) AddHeaders(w http.ResponseWriter) {
	w.Header().Set(ContentTypeHeader, ContentType)
	w.Header().Set(EtagHeader, fh.Watcher.CurrentSHA)
	w.Header().Set(CacheControlHeader, CacheHeaders)
}
