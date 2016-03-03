package server

import (
	"log"
	"net/http"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/vsco/dcdr/watcher"
	http_assert "github.com/vsco/goji-test/assert"
	builder "github.com/vsco/goji-test/builder"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/mutil"
)

func MockAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			watcher.AppendFeature("is_authorized", true, c)

			lw := mutil.WrapWriter(w)
			h.ServeHTTP(lw, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}

	return http.HandlerFunc(fn)
}

var features *watcher.FeatureMap

var cfg = &watcher.Config{
	ConfigPath:      "../config/decider_fixtures.json",
	FeatureEndpoint: "/decider.json",
}

func MockFeatures() *watcher.FeatureMap {
	if features != nil {
		return features
	}

	jsb, err := ioutil.ReadFile(cfg.ConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(jsb, &features)

	return features
}

func TestGetFeatures(t *testing.T) {
	srv := New(cfg, web.New()).Init()
	req := builder.WithMux(srv.Mux()).Get(cfg.FeatureEndpoint).Do()

	http_assert.Response(t, req.Response).
		IsOK().
		IsJSON().
		ContainsHeaderValue(handlers.CacheControlHeader, handlers.CacheHeaders).
		ContainsJSON(MockFeatures())
}

//func TestGetFeaturesCached(t *testing.T) {
//	srv := New(cfg, web.New()).Init()
//	etag := MockFeatures().Decider["current_sha"].(string)
//	hdrs := map[string]string{
//		handlers.IfNoneMatch: etag,
//	}
//
//	req := builder.WithMux(srv.Mux()).Get(cfg.FeatureEndpoint).Headers(hdrs).Do()
//
//	http_assert.Response(t, req.Response).
//		HasStatusCode(http.StatusNotModified).
//		ContainsHeaderValue("Cache-Control", "must-revalidate, public").
//		ContainsEtag(etag)
//}

func TestMiddleWareFeatures(t *testing.T) {
	srv := New(cfg, web.New()).Use(MockAuth).Init()

	req := builder.WithMux(srv.Mux()).Get(cfg.FeatureEndpoint).Do()

	assert.Equal(t, http.StatusUnauthorized, req.Response.Code)

	hdrs := map[string]string{
		"Authorization": "ABCDE",
	}

	req = builder.WithMux(srv.Mux()).Get(cfg.FeatureEndpoint).Headers(hdrs).Do()

	var fts *watcher.FeatureMap
	req.Response.UnmarshalBody(&fts)

	http_assert.Response(t, req.Response).IsJSON().IsOK()

	assert.True(t, fts.Decider["is_authorized"].(bool))
}
