package http

import (
	stats_api "github.com/fukata/golang-stats-api-handler"
	"goji.io/v3"
	"goji.io/v3/pat"

	api_v1 "github.com/oinume/lekcije/backend/proto_gen/proto/api/v1"
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
	mux.HandleFunc(pat.Get("/me/old"), s.getMeOld)
	mux.HandleFunc(pat.Get("/me/logout"), s.getMeLogout)
	mux.HandleFunc(pat.Get("/me/setting"), s.getMeSetting)
	mux.HandleFunc(pat.Get("/api/status"), s.getAPIStatus)
	mux.HandleFunc(pat.Get("/api/me/followingTeachers"), s.getAPIMeFollowingTeachers)
	mux.HandleFunc(pat.Get("/api/debug/envVar"), s.getAPIDebugEnvVar)
	mux.HandleFunc(pat.Get("/api/debug/httpHeader"), s.getAPIDebugHTTPHeader)
	mux.HandleFunc(pat.Post("/api/webhook/sendGrid"), s.postAPISendGridEventWebhook)
	mux.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)
}

func SetupTwirpServer(mux *goji.Mux, meServer api_v1.TwirpServer) {
	// TODO: Define UserServer etc in *server and setup mux inside it
	mux.Handle(pat.Post(api_v1.MePathPrefix+"*"), meServer)
}
