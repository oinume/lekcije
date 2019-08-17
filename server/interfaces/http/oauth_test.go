package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	r := require.New(t)
	s := newTestServer(t, nil)

	testCases := []struct {
		path        string
		handler     http.HandlerFunc
		code        int
		headerNames []string
	}{
		{
			path:    "/",
			handler: s.oauthGoogleHandler(),
			code:    http.StatusFound,
			headerNames: []string{
				"Set-Cookie",
				"Location",
			},
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", tc.path, nil)
		r.NoError(err)

		ctx := context_data.SetTrackingID(req.Context(), "a")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		tc.handler.ServeHTTP(w, req)

		a.Equal(tc.code, w.Code)
		for _, header := range tc.headerNames {
			if w.Header().Get(header) == "" {
				a.Failf("No header: %v", header)
			}
		}
	}
}
