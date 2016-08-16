package web

// TODO: Create package 'oauth'

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	google_auth2 "google.golang.org/api/oauth2/v2"
)

var _ = fmt.Print
var googleOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint:     google.Endpoint,
	RedirectURL:  "",
	Scopes: []string{
		"openid email",
		"openid profile",
	},
}

func OAuthGoogle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	state := util.RandomString(32)
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	c := getGoogleOAuthConfig(fmt.Sprintf("%s://%s/oauth/google/callback", scheme, r.Host))
	http.Redirect(w, r, c.AuthCodeURL(state), http.StatusFound)
}

func OAuthGoogleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := checkState(r); err != nil {
		InternalServerError(w, err)
		return
	}
	token, idToken, err := exchange(r)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	name, email, err := getNameAndEmail(token, idToken)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	db := model.MustDb(ctx)
	user := model.User{Name: name, Email: email}
	if err := db.FirstOrCreate(&user, model.User{Email: email}).Error; err != nil {
		InternalServerError(w, errors.Wrap(err, "Failed to get or create User"))
		return
	}

	// Create and save API Token
	apiToken := util.RandomString(64)
	userApiToken := model.UserApiToken{
		UserId: user.Id,
		Token:  apiToken,
	}
	if err := db.Create(&userApiToken).Error; err != nil {
		InternalServerError(w, errors.Wrap(err, "Failed to create UserApiToken"))
		return
	}
	cookie := &http.Cookie{
		Name:     ApiTokenCookieName,
		Value:    apiToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	//data := map[string]interface{}{
	//	"id":          user.Id,
	//	"name":        user.Name,
	//	"email":       user.Email,
	//	"accessToken": token.AccessToken,
	//	"idToken":     idToken,
	//}
	//if err := json.NewEncoder(w).Encode(data); err != nil {
	//	InternalServerError(w, errors.Errorf("Failed to encode JSON"))
	//	return
	//}

	http.Redirect(w, r, "/", http.StatusFound)
}

func checkState(r *http.Request) error {
	state := r.FormValue("state")
	oauthState, err := r.Cookie("oauthState")
	if err != nil {
		return errors.Wrap(err, "Failed to get cookie oauthState")
	}
	if state != oauthState.Value {
		return errors.Wrap(err, "state mismatch")
	}
	return nil
}

func exchange(r *http.Request) (*oauth2.Token, string, error) {
	code := r.FormValue("code")
	c := getGoogleOAuthConfig(fmt.Sprintf("http://%s/oauth/google/callback", r.Host)) // TODO: scheme
	token, err := c.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to exchange")
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.Errorf("Failed to get id_token")
	}
	return token, idToken, nil
}

func getNameAndEmail(token *oauth2.Token, idToken string) (string, string, error) {
	oauth2Client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))
	service, err := google_auth2.New(oauth2Client)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to create oauth2.Client")
	}

	userinfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get userinfo")
	}

	tokeninfo, err := service.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get tokeninfo")
	}

	return userinfo.Name, tokeninfo.Email, nil
}

func getGoogleOAuthConfig(redirectUrl string) oauth2.Config {
	c := googleOAuthConfig
	c.RedirectURL = redirectUrl
	return c
}
