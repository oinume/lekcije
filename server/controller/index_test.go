package controller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	r := require.New(t)
	err := os.Chdir("../..")
	r.NoError(err)

	server := httptest.NewServer(http.HandlerFunc(Index))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	r.NoError(err)
	ctx := context_data.SetTrackingID(req.Context(), "a")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	http.HandlerFunc(Index).ServeHTTP(rr, req)
	a.Equal(http.StatusOK, rr.Code)
	a.Contains(rr.Body.String(), "<title>lekcije - DMM英会話のお気に入り講師をフォローしよう</title>")
}
