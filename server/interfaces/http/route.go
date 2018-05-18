package http

import (
	"github.com/fukata/golang-stats-api-handler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"goji.io"
	"goji.io/pat"
)

func CreateMux(gatewayMux *runtime.ServeMux) *goji.Mux {
	s := NewServer()
	mux := goji.NewMux()
	mux.Use(SetTrackingID)
	mux.Use(AccessLogger)
	mux.Use(Redirecter)
	mux.Use(PanicHandler)
	mux.Use(NewRelic)
	mux.Use(SetDBAndRedis)
	mux.Use(SetLoggedInUser)
	mux.Use(LoginRequiredFilter)
	mux.Use(CORS)
	mux.Use(SetGRPCMetadata)
	mux.Use(SetGAMeasurementEventValues)

	mux.HandleFunc(pat.Get("/static/*"), s.static)
	mux.HandleFunc(pat.Get("/"), s.index)
	mux.HandleFunc(pat.Get("/signup"), s.signup)
	mux.HandleFunc(pat.Get("/oauth/google"), OAuthGoogle)
	mux.HandleFunc(pat.Get("/oauth/google/callback"), OAuthGoogleCallback)
	mux.HandleFunc(pat.Get("/robots.txt"), s.robotsTxt)
	mux.HandleFunc(pat.Get("/sitemap.xml"), s.sitemapXML)
	mux.HandleFunc(pat.Get("/terms"), s.terms)
	mux.HandleFunc(pat.Get("/me"), GetMe)
	mux.HandleFunc(pat.Post("/me/followingTeachers/create"), PostMeFollowingTeachersCreate)
	mux.HandleFunc(pat.Post("/me/followingTeachers/delete"), PostMeFollowingTeachersDelete)
	mux.HandleFunc(pat.Get("/me/logout"), GetMeLogout)
	mux.HandleFunc(pat.Get("/me/setting"), GetMeSetting)
	mux.HandleFunc(pat.Get("/api/status"), GetAPIStatus)
	mux.HandleFunc(pat.Get("/api/me/followingTeachers"), GetAPIMeFollowingTeachers)
	mux.HandleFunc(pat.Get("/api/debug/envVar"), GetAPIDebugEnvVar)
	mux.HandleFunc(pat.Get("/api/debug/httpHeader"), GetAPIDebugHTTPHeader)
	mux.HandleFunc(pat.Post("/api/webhook/sendGrid"), PostAPISendGridEventWebhook)
	mux.HandleFunc(pat.Post("/api/sendGrid/eventWebhook"), PostAPISendGridEventWebhook)
	mux.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)

	if gatewayMux != nil {
		// This path and path in the proto must be the same
		mux.Handle(pat.Get("/api/v1/*"), gatewayMux)
		mux.Handle(pat.Post("/api/v1/*"), gatewayMux)
	}

	return mux
}
