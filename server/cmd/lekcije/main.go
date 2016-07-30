package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oinume/lekcije/server/web"
	"github.com/oinume/lekcije/server/web/api"
	"github.com/oinume/lekcije/server/web/middleware"
	"goji.io"
	"goji.io/pat"
)

// TODO: move somewhere proper
var definedEnvs = map[string]string{
	"GOOGLE_CLIENT_ID":     "",
	"GOOGLE_CLIENT_SECRET": "",
	"DB_DSN":               "",
}

func init() {
	// Check env
	for key, _ := range definedEnvs {
		if value := os.Getenv(key); value != "" {
			definedEnvs[key] = value
		} else {
			log.Fatalf("Env '%v' must be defined.", key)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	mux := mux()
	fmt.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

func mux() *goji.Mux {
	mux := goji.NewMux()
	mux.UseC(middleware.SetDbToContext)
	mux.UseC(middleware.SetLoggedInUserToContext)
	mux.UseC(middleware.LoginRequiredFilter)

	mux.HandleFuncC(pat.Get("/"), web.Index)
	mux.HandleFuncC(pat.Get("/oauth/google"), web.OAuthGoogle)
	mux.HandleFuncC(pat.Get("/oauth/google/callback"), web.OAuthGoogleCallback)
	mux.HandleFuncC(pat.Post("/me/followingTeachers/create"), web.PostMeFollowingTeachersCreate)

	mux.HandleFuncC(pat.Get("/api/status"), api.GetStatus)
	mux.HandleFuncC(pat.Get("/api/me/followingTeachers"), api.GetMeFollowingTeachers)
	return mux
}
