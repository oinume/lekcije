package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	r := require.New(t)
	err := os.Chdir("../../..")
	r.NoError(err)

	s := NewServer()
	testCases := []struct {
		path     string
		handler  http.HandlerFunc
		code     int
		keywords []string
	}{
		{
			path:     "/",
			handler:  s.indexHandler(),
			code:     http.StatusOK,
			keywords: []string{"<title>lekcije - DMM英会話のお気に入り講師をフォローしよう</title>"},
		},
		{
			path:     "/robots.txt",
			handler:  s.robotsTxtHandler(),
			code:     http.StatusOK,
			keywords: []string{"Allow: /"},
		},
		{
			path:    "/sitemap.xml",
			handler: s.sitemapXMLHandler(),
			code:    http.StatusOK,
			keywords: []string{
				`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
				"/signup",
				"/terms",
			},
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", tc.path, nil)
		r.NoError(err)

		ctx := context_data.SetTrackingID(req.Context(), "a")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		tc.handler.ServeHTTP(rr, req)
		a.Equal(tc.code, rr.Code)
		for _, keyword := range tc.keywords {
			a.Contains(rr.Body.String(), keyword)
		}
	}
}
