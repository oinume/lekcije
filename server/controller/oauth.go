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

type oauthError int

const (
	oauthErrorUnknown oauthError = 1 + iota
	oauthErrorAccessDenied
)

func (e oauthError) Error() string {
	switch e {
	case oauthErrorUnknown:
		return "oauthError: unknown"
	case oauthErrorAccessDenied:
		return "oauthError: access denied"
	}
	return fmt.Sprintf("oauthError: invalid: %v", e)
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
		if err == oauthErrorAccessDenied {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
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
	if err == nil {
		go sendMeasurementEvent(r, eventCategoryAccount, "login", "via:google", 0)
	} else {
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
		go sendMeasurementEvent(r, eventCategoryAccount, "create", "via:google", 0)
	}

	userAPITokenService := model.NewUserAPITokenService(model.MustDB(ctx))
	userAPIToken, err := userAPITokenService.Create(user.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     APITokenCookieName,
		Value:    userAPIToken.Token,
		Path:     "/",
		Expires:  time.Now().Add(model.UserAPITokenExpiration),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/me", http.StatusFound)
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
	if e := r.FormValue("error"); e != "" {
		switch e {
		case "access_denied":
			return nil, "", oauthErrorAccessDenied
		default:
			return nil, "", oauthErrorUnknown
		}
	}
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
	c.RedirectURL = fmt.Sprintf("%s://%s/oauth/google/callback", config.WebURLScheme(r), r.Host)
	return c
}
