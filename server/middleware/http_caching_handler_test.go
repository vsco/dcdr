package middleware

import (
	"net/http"
	"testing"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/config"
	http_assert "github.com/vsco/goji-test/assert"
	"github.com/vsco/goji-test/builder"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func TestHTTPCaching(t *testing.T) {
	cfg := config.TestConfig()
	fm := models.EmptyFeatureMap()
	fm.Dcdr.Info.CurrentSHA = "current-sha"
	dcdr := client.New(cfg).SetFeatureMap(fm)
	mux := goji.DefaultMux

	mux.Get("/", func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Use(HTTPCachingHandler(dcdr))

	resp := builder.WithMux(mux).
		Get("/").
		Header(IfNoneMatchHeader, fm.Dcdr.Info.CurrentSHA).Do()

	http_assert.Response(t, resp.Response).
		HasStatusCode(http.StatusNotModified).
		ContainsHeaderValue(EtagHeader, fm.Dcdr.CurrentSHA()).
		ContainsHeaderValue(CacheControlHeader, CacheControl).
		ContainsHeaderValue(PragmaHeader, Pragma).
		ContainsHeaderValue(ExpiresHeader, Expires)

	resp = builder.WithMux(mux).
		Get("/").
		Header(IfNoneMatchHeader, "").Do()

	http_assert.Response(t, resp.Response).
		HasStatusCode(http.StatusOK).
		ContainsHeaderValue(EtagHeader, fm.Dcdr.CurrentSHA()).
		ContainsHeaderValue(CacheControlHeader, CacheControl).
		ContainsHeaderValue(PragmaHeader, Pragma).
		ContainsHeaderValue(ExpiresHeader, Expires)
}
