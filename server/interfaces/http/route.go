package http

import (
	"github.com/fukata/golang-stats-api-handler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/server/interfaces"
	"goji.io"
	"goji.io/pat"
)

func (s *server) CreateRoutes(gatewayMux *runtime.ServeMux, args *interfaces.ServerArgs) *goji.Mux {
	mux := goji.NewMux()
	mux.Use(SetTrackingID)
	mux.Use(AccessLogger)
	mux.Use(Redirecter)
	mux.Use(PanicHandler)
	mux.Use(NewRelic)
	mux.Use(SetDBAndRedis)
	mux.Use(setLoggedInUser(args.DB))
	mux.Use(loginRequiredFilter(args.DB))
	mux.Use(CORS)
	mux.Use(SetGRPCMetadata)
	mux.Use(SetGAMeasurementEventValues)

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

	if gatewayMux != nil {
		// This path and path in the proto must be the same
		mux.Handle(pat.Get("/api/v1/*"), gatewayMux)
		mux.Handle(pat.Post("/api/v1/*"), gatewayMux)
	}

	return mux
}
