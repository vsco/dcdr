package server

import (
	"net/http"
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	http_assert "github.com/vsco/goji-test/assert"
	"github.com/vsco/goji-test/builder"
	"github.com/zenazn/goji"
)

var fm = models.EmptyFeatureMap()
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
		IsJSON()

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
		{"a", []string{"a"}},
	}

	for _, tc := range cases {
		r.Header.Set(handlers.DcdrScopesHeader, tc.Input)
		assert.Equal(t, tc.Expected, handlers.GetScopes(r), tc.Input)
	}
}
