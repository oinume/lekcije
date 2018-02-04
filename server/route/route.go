package route

import (
	"github.com/fukata/golang-stats-api-handler"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/middleware"
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

	routes.HandleFunc(pat.Get("/static/*"), controller.Static)
	routes.HandleFunc(pat.Get("/"), controller.Index)
	routes.HandleFunc(pat.Get("/signup"), controller.Signup)
	routes.HandleFunc(pat.Get("/oauth/google"), controller.OAuthGoogle)
	routes.HandleFunc(pat.Get("/oauth/google/callback"), controller.OAuthGoogleCallback)
	routes.HandleFunc(pat.Get("/robots.txt"), controller.RobotsTxt)
	routes.HandleFunc(pat.Get("/sitemap.xml"), controller.SitemapXML)
	routes.HandleFunc(pat.Get("/terms"), controller.Terms)
	routes.HandleFunc(pat.Get("/me"), controller.GetMe)
	routes.HandleFunc(pat.Post("/me/followingTeachers/create"), controller.PostMeFollowingTeachersCreate)
	routes.HandleFunc(pat.Post("/me/followingTeachers/delete"), controller.PostMeFollowingTeachersDelete)
	routes.HandleFunc(pat.Get("/me/logout"), controller.GetMeLogout)
	routes.HandleFunc(pat.Get("/me/setting"), controller.GetMeSetting)
	routes.HandleFunc(pat.Post("/me/setting/update"), controller.PostMeSettingUpdate)
	routes.HandleFunc(pat.Get("/api/status"), controller.GetAPIStatus)
	routes.HandleFunc(pat.Get("/api/me/followingTeachers"), controller.GetAPIMeFollowingTeachers)
	routes.HandleFunc(pat.Get("/api/debug/envVar"), controller.GetAPIDebugEnvVar)
	routes.HandleFunc(pat.Get("/api/debug/httpHeader"), controller.GetAPIDebugHTTPHeader)
	routes.HandleFunc(pat.Post("/api/webhook/sendGrid"), controller.PostAPISendGridEventWebhook)
	routes.HandleFunc(pat.Post("/api/sendGrid/eventWebhook"), controller.PostAPISendGridEventWebhook)
	routes.HandleFunc(pat.Get("/api/stats"), stats_api.Handler)

	if gatewayMux != nil {
		// This path and path in the proto must be the same
		routes.Handle(pat.Get("/api/v1/echo"), gatewayMux)
		routes.Handle(pat.Get("/api/v2/echo"), gatewayMux)
		//routes.Handle(pat.Get("/api/v1/setting/email"), gatewayMux)
		routes.Handle(pat.Get("/api/v1/me"), gatewayMux)
		routes.Handle(pat.Get("/api/v1/me/email"), gatewayMux)
		routes.Handle(pat.Post("/api/v1/me/email"), gatewayMux)
		routes.Handle(pat.Post("/api/v1/me/notificationTimeSpan"), gatewayMux)
	}

	return routes
}
