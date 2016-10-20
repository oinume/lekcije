package controller

// TODO: Create package 'oauth'

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
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

func OAuthGoogle(w http.ResponseWriter, r *http.Request) {
	state := util.RandomString(32)
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	c := getGoogleOAuthConfig(r)
	http.Redirect(w, r, c.AuthCodeURL(state), http.StatusFound)
}

func OAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := checkState(r); err != nil {
		InternalServerError(w, err)
		return
	}
	token, idToken, err := exchange(r)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	googleID, name, email, err := getGoogleUserInfo(token, idToken)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	db := model.MustDB(ctx)
	userService := model.NewUserService(db)
	user, err := userService.FindByGoogleID(googleID)
	if err != nil {
		if _, notFound := err.(*errors.NotFound); !notFound {
			InternalServerError(w, err)
			return
		}
		// Couldn't find user for the googleID, so create a new user
		errTx := model.GORMTransaction(db, "OAuthGoogleCallback", func(tx *gorm.DB) error {
			var errCreate error
			user, _, errCreate = userService.CreateWithGoogle(name, email, googleID)
			return errCreate
		})
		if errTx != nil {
			InternalServerError(w, errTx)
			return
		}
	}

	userApiTokenService := model.NewUserApiTokenService(model.MustDB(ctx))
	userApiToken, err := userApiTokenService.Create(user.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     ApiTokenCookieName,
		Value:    userApiToken.Token,
		Path:     "/",
		Expires:  time.Now().Add(model.UserApiTokenExpiration),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	//data := map[string]interface{}{
	//	"id":          user.ID,
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
		return errors.InternalWrapf(err, "Failed to get cookie oauthState")
	}
	if state != oauthState.Value {
		return errors.InternalWrapf(err, "state mismatch")
	}
	return nil
}

func exchange(r *http.Request) (*oauth2.Token, string, error) {
	code := r.FormValue("code")
	c := getGoogleOAuthConfig(r)
	token, err := c.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, "", errors.InternalWrapf(err, "Failed to exchange")
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.Internalf("Failed to get id_token")
	}
	return token, idToken, nil
}

// Returns userId, name, email, error
func getGoogleUserInfo(token *oauth2.Token, idToken string) (string, string, string, error) {
	oauth2Client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))
	service, err := google_auth2.New(oauth2Client)
	if err != nil {
		// TODO: quit using errors.Wrap
		return "", "", "", errors.InternalWrapf(err, "Failed to create oauth2.Client")
	}

	userinfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", "", "", errors.InternalWrapf(err, "Failed to get userinfo")
	}

	tokeninfo, err := service.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return "", "", "", errors.InternalWrapf(err, "Failed to get tokeninfo")
	}

	return tokeninfo.UserId, userinfo.Name, tokeninfo.Email, nil
}

func getGoogleOAuthConfig(r *http.Request) oauth2.Config {
	c := googleOAuthConfig
	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	if config.IsProductionEnv() {
		scheme = "https"
	}
	c.RedirectURL = fmt.Sprintf("%s://%s/oauth/google/callback", scheme, r.Host)
	return c
}
