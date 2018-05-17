package route

import (
	"github.com/fukata/golang-stats-api-handler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/server/interfaces/http"
	"github.com/oinume/lekcije/server/interfaces/http/middleware"
	"goji.io"
	"goji.io/pat"
)

func Create(gatewayMux *runtime.ServeMux) *goji.Mux {
	routes := goji.NewMux()
	routes.Use(middleware.SetTrackingID)
	routes.Use(middleware.AccessLogger)
	routes.Use(middleware.Redirecter)
	routes.Use(middleware.PanicHandler)
	routes.Use(middleware.NewRelic)
	routes.Use(middleware.SetDBAndRedis)
	routes.Use(middleware.SetLoggedInUser)
	routes.Use(middleware.LoginRequiredFilter)
	routes.Use(middleware.CORS)
	routes.Use(middleware.SetGRPCMetadata)
	routes.Use(middleware.SetGAMeasurementEventValues)

	routes.HandleFunc(pat.Get("/static/*"), http.Static)
	routes.HandleFunc(pat.Get("/"), http.Index)
	routes.HandleFunc(pat.Get("/signup"), http.Signup)
	routes.HandleFunc(pat.Get("/oauth/google"), http.OAuthGoogle)
	routes.HandleFunc(pat.Get("/oauth/google/callback"), http.OAuthGoogleCallback)
	routes.HandleFunc(pat.Get("/robots.txt"), http.RobotsTxt)
	routes.HandleFunc(pat.Get("/sitemap.xml"), http.SitemapXML)
	routes.HandleFunc(pat.Get("/terms"), http.Terms)
	routes.HandleFunc(pat.Get("/me"), http.GetMe)
	routes.HandleFunc(pat.Post("/me/followingTeachers/create"), http.PostMeFollowingTeachersCreate)
	routes.HandleFunc(pat.Post("/me/followingTeachers/delete"), http.PostMeFollowingTeachersDelete)
	routes.HandleFunc(pat.Get("/me/logout"), http.GetMeLogout)
	routes.HandleFunc(pat.Get("/me/setting"), http.GetMeSetting)
	routes.HandleFunc(pat.Get("/api/status"), http.GetAPIStatus)
	routes.HandleFunc(pat.Get("/api/me/followingTeachers"), http.GetAPIMeFollowingTeachers)
	routes.HandleFunc(pat.Get("/api/debug/envVar"), http.GetAPIDebugEnvVar)
	routes.HandleFunc(pat.Get("/api/debug/httpHeader"), http.GetAPIDebugHTTPHeader)
	routes.HandleFunc(pat.Post("/api/webhook/sendGrid"), http.PostAPISendGridEventWebhook)
	routes.HandleFunc(pat.Post("/api/sendGrid/eventWebhook"), http.PostAPISendGridEventWebhook)
	routes.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)

	if gatewayMux != nil {
		// This path and path in the proto must be the same
		routes.Handle(pat.Get("/api/v1/*"), gatewayMux)
		routes.Handle(pat.Post("/api/v1/*"), gatewayMux)
	}

	return routes
}
