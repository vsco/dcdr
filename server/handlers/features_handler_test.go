package handlers

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func R() *http.Request {
	r, _ := http.NewRequest("GET", "http://:9191/", bytes.NewReader([]byte{}))
	return r
}

func TestGetScopes(t *testing.T) {
	r := R()
	r.Header.Set(DcdrScopesHeader, " a,b, c")

	assert.Equal(t, []string{"a", "b", "c"}, GetScopes(r))
}

func TestSetScopes(t *testing.T) {
	r := R()
	scopes := []string{"d", "e"}
	SetScopes(r, scopes)

	assert.Equal(t, scopes, GetScopes(r))
}

func TestAppendScope(t *testing.T) {
	r := R()
	scope := "d"
	r.Header.Set(DcdrScopesHeader, "a/b/c")
	AppendScope(r, scope)

	assert.Equal(t, []string{"a/b/c", scope}, GetScopes(r))
}
