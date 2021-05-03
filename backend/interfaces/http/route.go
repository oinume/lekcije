package http

import (
	stats_api "github.com/fukata/golang-stats-api-handler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
	"goji.io/v3"
	"goji.io/v3/pat"
)

func (s *server) CreateRoutes(gatewayMux *runtime.ServeMux) *goji.Mux {
	mux := goji.NewMux()
	mux.Use(setTrackingID)
	mux.Use(accessLogger(s.accessLogger))
	mux.Use(redirecter)
	mux.Use(panicHandler)
	mux.Use(setLoggedInUser(s.db))
	mux.Use(loginRequiredFilter(s.db, s.appLogger))
	mux.Use(setCORS)
	mux.Use(setGRPCMetadata)
	mux.Use(setGAMeasurementEventValues)

	mux.HandleFunc(pat.Get("/static/*"), s.static)
	mux.HandleFunc(pat.Get("/"), s.index)
	mux.HandleFunc(pat.Get("/signup"), s.signup)
	mux.HandleFunc(pat.Get("/oauth/google"), s.oauthGoogle)
	mux.HandleFunc(pat.Get("/oauth/google/callback"), s.oauthGoogleCallback)
	mux.HandleFunc(pat.Get("/robots.txt"), s.robotsTxt)
	mux.HandleFunc(pat.Get("/sitemap.xml"), s.sitemapXML)
	mux.HandleFunc(pat.Get("/terms"), s.terms)
	mux.HandleFunc(pat.Get("/me"), s.getMe)
	mux.HandleFunc(pat.Post("/me/followingTeachers/create"), s.postMeFollowingTeachersCreate)
	mux.HandleFunc(pat.Post("/me/followingTeachers/delete"), s.postMeFollowingTeachersDelete)
	mux.HandleFunc(pat.Get("/me/logout"), s.getMeLogout)
	mux.HandleFunc(pat.Get("/me/setting"), s.getMeSetting)
	mux.HandleFunc(pat.Get("/api/status"), s.getAPIStatus)
	mux.HandleFunc(pat.Get("/api/me/followingTeachers"), s.getAPIMeFollowingTeachers)
	mux.HandleFunc(pat.Get("/api/debug/envVar"), s.getAPIDebugEnvVar)
	mux.HandleFunc(pat.Get("/api/debug/httpHeader"), s.getAPIDebugHTTPHeader)
	mux.HandleFunc(pat.Post("/api/webhook/sendGrid"), s.postAPISendGridEventWebhook)
	mux.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)

	userService := NewUserService()
	mux.Handle(pat.Post(api_v1.UserPathPrefix+"*"), api_v1.NewUserServer(userService))

	if gatewayMux != nil {
		// This path and path in the proto must be the same
		mux.Handle(pat.Get("/api/v1/*"), gatewayMux)
		mux.Handle(pat.Post("/api/v1/*"), gatewayMux)
	}

	return mux
}
