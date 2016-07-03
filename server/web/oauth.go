package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	google_auth2 "google.golang.org/api/oauth2/v2"

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
	state := randomString(32)
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, googleOAuthConfig.AuthCodeURL(state), http.StatusFound)
}

func OAuthGoogleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	checkState(w, r)
	token, idToken, err := exchange(r)
	if err != nil {
		internalServerError(w, err.Error())
	}
	name, email, err := getNameAndEmail(token, idToken)
	if err != nil {
		internalServerError(w, err.Error())
	}
	db, err := model.Open()
	if err != nil {
		internalServerError(w, fmt.Sprintf("Failed to connect db: %v", err))
	}

	user := model.User{Name: name, Email: email}
	if err := db.FirstOrCreate(&user, model.User{Email: email}).Error; err != nil {
		internalServerError(w, fmt.Sprintf("Failed to access user: %v", err))
	}

	data := map[string]interface{}{
		"id":          user.Id,
		"name":        user.Name,
		"email":       user.Email,
		"accessToken": token.AccessToken,
		"idToken":     idToken,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON"), http.StatusInternalServerError)
		return
	}
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func checkState(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	oauthState, err := r.Cookie("oauthState")
	if err != nil {
		internalServerError(w, fmt.Sprintf("Failed to get cookie oauthState: %s", err))
		return
	}
	if state != oauthState.Value {
		internalServerError(w, "state mismatch")
		return
	}
}

func exchange(r *http.Request) (*oauth2.Token, string, error) {
	code := r.FormValue("code")
	token, err := googleOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, "", err
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", fmt.Errorf("Failed to get id_token")
	}
	return token, idToken, nil
}

func getNameAndEmail(token *oauth2.Token, idToken string) (string, string, error) {
	oauth2Client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))
	service, err := google_auth2.New(oauth2Client)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create oauth2.Client: %v", err)
	}

	userinfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", "", err
		//internalServerError(w, fmt.Sprintf("Failed to get userinfo: %v", err))
	}

	tokeninfo, err := service.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return "", "", err
	}

	return userinfo.Name, tokeninfo.Email, nil
}
