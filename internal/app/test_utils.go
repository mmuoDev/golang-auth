package app

import (
	"net/http"
	"net/http/httptest"
)

func NewTestServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() { ts.Close() }
}
