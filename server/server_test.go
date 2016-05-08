package server

import (
	"net/http"
	"testing"

	"bytes"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/vsco/dcdr/server/middleware"
	http_assert "github.com/vsco/http-test/assert"
	"github.com/vsco/http-test/builder"
)

var fm = models.EmptyFeatureMap()
var cfg = config.TestConfig()
var cl, _ = client.New(cfg)

func mockServer() *Server {
	cl.SetFeatureMap(fm)
	return New(cfg, cl)
}

func TestGetFeatures(t *testing.T) {
	srv := mockServer()
	resp := builder.WithMux(srv).Get(srv.config.Server.Endpoint).Do()
	http_assert.Response(t, resp.Response).
		IsOK().
		IsJSON()

	var m models.FeatureMap
	err := resp.Response.UnmarshalBody(&m)

	assert.NoError(t, err)
	assert.Equal(t, cl.ScopedMap(), &m)
}

func TestScopeHeader(t *testing.T) {
	srv := mockServer()
	resp := builder.WithMux(srv).
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

func TestGetScopes(t *testing.T) {
	rd := bytes.NewReader([]byte{})
	r, err := http.NewRequest("GET", "/", rd)
	assert.NoError(t, err)

	cases := []struct {
		Input    string
		Expected []string
	}{
		{" scope1, scope2 ", []string{"scope1", "scope2"}},
		{"scope1,scope2", []string{"scope1", "scope2"}},
		{" a/b/c, d", []string{"a/b/c", "d"}},
		// Prevent unbounded scopes
		{"1,2,3,4,5,6,7,8,9", []string{"1", "2", "3", "4", "5", "6", "7", "8"}},
		{"a,a,a", []string{"a"}},
		{"a", []string{"a"}},
	}

	for _, tc := range cases {
		r.Header.Set(handlers.DcdrScopesHeader, tc.Input)
		assert.Equal(t, tc.Expected, handlers.GetScopes(r), tc.Input)
	}
}

func TestHTTPCaching(t *testing.T) {
	srv := mockServer()
	ts := time.Now().Unix()
	fm := models.EmptyFeatureMap()
	fm.Dcdr.Info.CurrentSHA = "current-sha"
	fm.Dcdr.Info.LastModifiedDate = ts
	srv.Client.SetFeatureMap(fm)

	resp := builder.WithMux(srv).
		Get(srv.config.Server.Endpoint).
		Header(middleware.IfNoneMatchHeader, fm.Dcdr.Info.CurrentSHA).Do()

	http_assert.Response(t, resp.Response).
		HasStatusCode(http.StatusNotModified).
		ContainsHeaderValue(middleware.EtagHeader, fm.Dcdr.CurrentSHA()).
		ContainsHeaderValue(middleware.LastModifiedHeader, time.Unix(ts, 0).Format(time.RFC1123)).
		ContainsHeaderValue(middleware.CacheControlHeader, middleware.CacheControl).
		ContainsHeaderValue(middleware.PragmaHeader, middleware.Pragma).
		ContainsHeaderValue(middleware.ExpiresHeader, middleware.Expires)

	resp = builder.WithMux(srv).
		Get(srv.config.Server.Endpoint).
		Header(middleware.IfNoneMatchHeader, "").Do()

	http_assert.Response(t, resp.Response).
		HasStatusCode(http.StatusOK).
		ContainsHeaderValue(middleware.EtagHeader, fm.Dcdr.CurrentSHA()).
		ContainsHeaderValue(middleware.CacheControlHeader, middleware.CacheControl).
		ContainsHeaderValue(middleware.PragmaHeader, middleware.Pragma).
		ContainsHeaderValue(middleware.ExpiresHeader, middleware.Expires)
}
