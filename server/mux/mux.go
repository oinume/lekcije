package mux

import (
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/middleware"
	"goji.io"
	"goji.io/pat"
)

func Create() *goji.Mux {
	mux := goji.NewMux()
	mux.Use(middleware.AccessLogger)
	mux.Use(middleware.Redirecter)
	mux.Use(middleware.PanicHandler)
	mux.Use(middleware.NewRelic)
	mux.Use(middleware.SetDBAndRedisToContext)
	mux.Use(middleware.SetLoggedInUserToContext)
	mux.Use(middleware.LoginRequiredFilter)
	mux.Use(middleware.CORS)

	mux.HandleFunc(pat.Get("/static/*"), controller.Static)
	mux.HandleFunc(pat.Get("/"), controller.Index)
	mux.HandleFunc(pat.Get("/logout"), controller.Logout)
	mux.HandleFunc(pat.Get("/oauth/google"), controller.OAuthGoogle)
	mux.HandleFunc(pat.Get("/oauth/google/callback"), controller.OAuthGoogleCallback)
	mux.HandleFunc(pat.Post("/me/followingTeachers/create"), controller.PostMeFollowingTeachersCreate)
	mux.HandleFunc(pat.Post("/me/followingTeachers/delete"), controller.PostMeFollowingTeachersDelete)
	mux.HandleFunc(pat.Get("/me/setting"), controller.GetMeSetting)
	mux.HandleFunc(pat.Post("/me/setting/update"), controller.PostMeSettingUpdate)
	mux.HandleFunc(pat.Get("/api/status"), controller.GetAPIStatus)
	mux.HandleFunc(pat.Get("/api/me/followingTeachers"), controller.GetAPIMeFollowingTeachers)

	return mux
}
