package route

import (
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/middleware"
	"goji.io"
	"goji.io/pat"
)

func Create() *goji.Mux {
	routes := goji.NewMux()
	routes.Use(middleware.SetTrackingID)
	routes.Use(middleware.AccessLogger)
	routes.Use(middleware.Redirecter)
	routes.Use(middleware.PanicHandler)
	routes.Use(middleware.NewRelic)
	routes.Use(middleware.SetDBAndRedisToContext)
	routes.Use(middleware.SetLoggedInUserToContext)
	routes.Use(middleware.LoginRequiredFilter)
	routes.Use(middleware.CORS)

	routes.HandleFunc(pat.Get("/static/*"), controller.Static)
	routes.HandleFunc(pat.Get("/"), controller.Index)
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

	return routes
}
