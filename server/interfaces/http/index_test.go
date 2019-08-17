package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	r := require.New(t)

	s := newTestServer(t, nil)
	testCases := []struct {
		path     string
		handler  http.HandlerFunc
		code     int
		keywords []string
	}{
		{
			path:    "/",
			handler: s.indexHandler(),
			code:    http.StatusOK,
			keywords: []string{
				`<title>lekcije - DMM英会話のお気に入り講師をフォローしよう</title>`,
			},
		},
		{
			path:     "/robots.txt",
			handler:  s.robotsTxtHandler(),
			code:     http.StatusOK,
			keywords: []string{`Allow: /`},
		},
		{
			path:    "/signup",
			handler: s.signupHandler(),
			code:    http.StatusOK,
			keywords: []string{
				`<title>新規登録 | lekcije</title>`,
			},
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
		{
			path:    "/terms",
			handler: s.termsHandler(),
			code:    http.StatusOK,
			keywords: []string{
				`<title>利用規約 | lekcije</title>`,
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
		for _, keyword := range tc.keywords {
			a.Contains(w.Body.String(), keyword)
		}
	}
}
