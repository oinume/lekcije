package mux

import (
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/api"
	"github.com/oinume/lekcije/server/controller/middleware"
	"goji.io"
	"goji.io/pat"
)

func Create() *goji.Mux {
	mux := goji.NewMux()
	mux.UseC(middleware.AccessLogger)
	mux.UseC(middleware.SetDbToContext)
	mux.UseC(middleware.SetLoggedInUserToContext)
	mux.UseC(middleware.LoginRequiredFilter)
	mux.UseC(middleware.CORS)

	mux.HandleFunc(pat.Get("/static/*"), controller.Static)
	mux.HandleFuncC(pat.Get("/"), controller.Index)
	mux.HandleFuncC(pat.Get("/logout"), controller.Logout)
	mux.HandleFuncC(pat.Get("/oauth/google"), controller.OAuthGoogle)
	mux.HandleFuncC(pat.Get("/oauth/google/callback"), controller.OAuthGoogleCallback)
	mux.HandleFuncC(pat.Post("/me/followingTeachers/create"), controller.PostMeFollowingTeachersCreate)
	mux.HandleFuncC(pat.Post("/me/followingTeachers/delete"), controller.PostMeFollowingTeachersDelete)
	mux.HandleFuncC(pat.Get("/me/setting"), controller.GetMeSetting)
	mux.HandleFuncC(pat.Get("/me/setting/update"), controller.PostMeSettingUpdate)
	mux.HandleFuncC(pat.Get("/api/status"), api.GetStatus)
	mux.HandleFuncC(pat.Get("/api/me/followingTeachers"), api.GetMeFollowingTeachers)

	return mux
}
