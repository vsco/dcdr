package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	http_assert "github.com/vsco/goji-test/assert"
	"github.com/vsco/goji-test/builder"
	"github.com/zenazn/goji"
)

var fm = &models.FeatureMap{
	Dcdr: models.Root{
		Info: models.Info{
			CurrentSha: "abcde",
		},
		Features: models.Features{
			"default": map[string]interface{}{
				"bool":  true,
				"float": 0.5,
			},
			"scope": map[string]interface{}{
				"scope-bool":  true,
				"scope-float": 0.5,
			},
		},
	},
}

var cfg = config.TestConfig()
var cl = client.New(cfg).SetFeatureMap(fm)

func Server() *server {
	return New(cfg, goji.DefaultMux, cl).BindMux()
}

func TestGetFeatures(t *testing.T) {
	srv := Server()
	resp := builder.WithMux(srv.mux).Get(srv.config.Server.Endpoint).Do()
	http_assert.Response(t, resp.Response).
		IsOK().
		IsJSON().
		ContainsHeaderValue(handlers.EtagHeader, fm.Dcdr.CurrentSha()).
		ContainsHeaderValue(handlers.CacheControlHeader, handlers.CacheControl).
		ContainsHeaderValue(handlers.PragmaHeader, handlers.Pragma).
		ContainsHeaderValue(handlers.ExpiresHeader, handlers.Expires)

	var m models.FeatureMap
	err := resp.Response.UnmarshalBody(&m)

	assert.NoError(t, err)
	assert.Equal(t, cl.ScopedMap(), &m)
}

func TestScopeHeader(t *testing.T) {
	srv := Server()
	resp := builder.WithMux(srv.mux).
		Get(srv.config.Server.Endpoint).
		Header(handlers.DcdrScopesHeader, "scope, scope2").Do()

	http_assert.Response(t, resp.Response).
		IsOK().
		IsJSON().
		ContainsHeaderValue(handlers.DcdrScopesHeader, "scope, scope2")

	var m models.FeatureMap
	err := resp.Response.UnmarshalBody(&m)

	assert.NoError(t, err)
	assert.Equal(t, cl.WithScopes("scope").ScopedMap(), &m)
}

func TestHTTPCaching(t *testing.T) {
	srv := Server()
	resp := builder.WithMux(srv.mux).
		Get(srv.config.Server.Endpoint).
		Header(handlers.IfNoneMatchHeader, fm.Dcdr.CurrentSha()).Do()
	http_assert.Response(t, resp.Response).
		HasStatusCode(http.StatusNotModified)
}
