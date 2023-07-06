package http

import (
	stats_api "github.com/fukata/golang-stats-api-handler"
	"goji.io/v3"
	"goji.io/v3/pat"
)

func (s *server) Setup(mux *goji.Mux) {
	mux.Use(setTrackingID)
	mux.Use(accessLogger(s.accessLogger))
	mux.Use(redirecter)
	mux.Use(panicHandler(s.errorRecorder))
	mux.Use(setLoggedInUser(s.db))
	mux.Use(loginRequiredFilter(s.db, s.appLogger, s.errorRecorder))
	mux.Use(setCORS)
	mux.Use(setGAMeasurementEventValues)
	mux.Use(setAuthorizationContext)

	mux.HandleFunc(pat.Get("/static/*"), s.static)
	mux.HandleFunc(pat.Get("/"), s.index)
	mux.HandleFunc(pat.Get("/signup"), s.signup)
	mux.HandleFunc(pat.Get("/robots.txt"), s.robotsTxt)
	mux.HandleFunc(pat.Get("/sitemap.xml"), s.sitemapXML)
	mux.HandleFunc(pat.Get("/terms"), s.terms)
	mux.HandleFunc(pat.Get("/me"), s.getMe)
	mux.HandleFunc(pat.Get("/me/logout"), s.getMeLogout)
	mux.HandleFunc(pat.Get("/me/setting"), s.getMeSetting)
	mux.HandleFunc(pat.Get("/api/status"), s.getAPIStatus)
	mux.HandleFunc(pat.Get("/api/debug/envVar"), s.getAPIDebugEnvVar)
	mux.HandleFunc(pat.Get("/api/debug/httpHeader"), s.getAPIDebugHTTPHeader)
	mux.HandleFunc(pat.Post("/api/webhook/sendGrid"), s.postAPISendGridEventWebhook)
	mux.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)
}
