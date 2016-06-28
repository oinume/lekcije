package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/oinume/lekcije/server/model"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint:     google.Endpoint,
	RedirectURL:  fmt.Sprintf("http://localhost:%d/oauth/google/callback", 4000),
	Scopes: []string{
		"openid email",
		"openid profile",
	},
}

func OAuthGoogle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: Give state
	http.Redirect(w, r, googleOAuthConfig.AuthCodeURL(""), http.StatusFound)
}

func OAuthGoogleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: state check
	code := r.FormValue("code")
	token, err := googleOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Exchange error: %s", err), http.StatusInternalServerError)
		return
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, fmt.Sprintf("Failed to get id_token"), http.StatusInternalServerError)
		return
	}

	db, err := model.Open()
	user := model.User{Name: "oinume", Email: "oinume@gmail.com"}
	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}

	data := map[string]string{
		"accessToken": token.AccessToken,
		"expiry":      token.Expiry.Format(time.RFC3339),
		"idToken":     idToken,
		"tokenType":   token.TokenType,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON"), http.StatusInternalServerError)
		return
	}
}
