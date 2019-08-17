package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/lekcije/server/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessLogger(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	req, err := http.NewRequest("GET", "/", nil)
	r.NoError(err)
	w := httptest.NewRecorder()
	b := new(bytes.Buffer)
	middleware := accessLogger(logger.NewAccessLogger(b))
	middleware(http.HandlerFunc(dummyHandler)).ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)
	a.Contains(b.String(), `"access"`)
}

func TestSeTrackingID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	req, err := http.NewRequest("GET", "/", nil)
	r.NoError(err)
	w := httptest.NewRecorder()
	middleware := setTrackingID(http.HandlerFunc(dummyHandler))
	middleware.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)
	a.Contains(w.Header().Get("Set-Cookie"), TrackingIDCookieName)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
