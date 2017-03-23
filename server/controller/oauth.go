package controller

// TODO: Create package 'oauth'

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	facebook_api "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
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

var facebookOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
	ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
	Endpoint:     facebook.Endpoint,
	RedirectURL:  "",
	Scopes: []string{
		"email",
		"public_profile",
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
	token, idToken, err := exchangeGoogle(r)
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

	db := context_data.MustDB(ctx)
	userService := model.NewUserService(db)
	user, err := userService.FindByGoogleID(googleID)
	if err == nil {
		go sendMeasurementEvent(r, eventCategoryUser, "login", fmt.Sprint(user.ID), 0, user.ID)
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
		go sendMeasurementEvent(r, eventCategoryUser, "create", fmt.Sprint(user.ID), 0, user.ID)
	}

	userAPITokenService := model.NewUserAPITokenService(context_data.MustDB(ctx))
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

func OAuthFacebook(w http.ResponseWriter, r *http.Request) {
	state := util.RandomString(32)
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	c := getFacebookOAuthConfig(r)
	http.Redirect(w, r, c.AuthCodeURL(state), http.StatusFound)
}

func OAuthFacebookCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := checkState(r); err != nil {
		InternalServerError(w, err)
		return
	}
	token, _, err := exchangeFacebook(r)
	if err != nil {
		if err == oauthErrorAccessDenied {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		InternalServerError(w, err)
		return
	}

	facebookID, name, email, err := getFacebookUserInfo(token)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	db := context_data.MustDB(ctx)
	userService := model.NewUserService(db)
	user, err := userService.FindByFacebookID(facebookID)
	if err == nil {
		go sendMeasurementEvent(r, eventCategoryUser, "login", fmt.Sprint(user.ID), 0, user.ID)
	} else {
		if _, notFound := err.(*errors.NotFound); !notFound {
			InternalServerError(w, err)
			return
		}
		// Couldn't find user for the facebookID, so create a new user
		errTx := model.GORMTransaction(db, "OAuthFacebookCallback", func(tx *gorm.DB) error {
			var errCreate error
			user, _, errCreate = userService.CreateWithFacebook(name, email, facebookID)
			return errCreate
		})
		if errTx != nil {
			InternalServerError(w, errTx)
			return
		}
		go sendMeasurementEvent(r, eventCategoryUser, "create", fmt.Sprint(user.ID), 0, user.ID)
	}

	cookie, err := createUserAPIToken(context_data.MustDB(ctx), user.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/me", http.StatusFound)
}

func createUserAPIToken(db *gorm.DB, userID uint32) (*http.Cookie, error) {
	userAPITokenService := model.NewUserAPITokenService(db)
	userAPIToken, err := userAPITokenService.Create(userID)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     APITokenCookieName,
		Value:    userAPIToken.Token,
		Path:     "/",
		Expires:  time.Now().Add(model.UserAPITokenExpiration),
		HttpOnly: false,
	}, nil
}

func checkState(r *http.Request) error {
	state := r.FormValue("state")
	oauthState, err := r.Cookie("oauthState")
	if err != nil {
		return errors.InternalWrapf(
			err, "Failed to get cookie oauthState: userAgent=%v, remoteAddr=%v",
			r.UserAgent(), GetRemoteAddress(r),
		)
	}
	if state != oauthState.Value {
		return errors.InternalWrapf(err, "state mismatch")
	}
	return nil
}

func exchangeGoogle(r *http.Request) (*oauth2.Token, string, error) {
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
	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", errors.InternalWrapf(err, "Failed to exchange")
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.Internalf("Failed to get id_token")
	}
	return token, idToken, nil
}

func exchangeFacebook(r *http.Request) (*oauth2.Token, string, error) {
	if e := r.FormValue("error"); e != "" {
		switch e {
		case "access_denied":
			return nil, "", oauthErrorAccessDenied
		default:
			return nil, "", oauthErrorUnknown
		}
	}
	code := r.FormValue("code")
	c := getFacebookOAuthConfig(r)
	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", errors.InternalWrapf(err, "Failed to exchange")
	}
	//idToken, ok := token.Extra("id_token").(string)
	//if !ok {
	//	return nil, "", errors.Internalf("Failed to get id_token")
	//}
	return token, "", nil
}

// Returns userId, name, email, error
func getGoogleUserInfo(token *oauth2.Token, idToken string) (string, string, string, error) {
	oauth2Client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	oauth2Client.Timeout = 5 * time.Second
	service, err := google_auth2.New(oauth2Client)
	if err != nil {
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

// TODO: return model.User
func getFacebookUserInfo(token *oauth2.Token) (string, string, string, error) {
	oauth2Client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	oauth2Client.Timeout = 5 * time.Second
	fb := facebook_api.New(os.Getenv("FACEBOOK_CLIENT_ID"), os.Getenv("FACEBOOK_CLIENT_SECRET"))
	session := fb.Session(token.AccessToken)
	session.HttpClient = oauth2Client

	result, err := session.Get("/me", facebook_api.Params{
		"fields": "id,email,name",
	})
	if err != nil {
		return "", "", "", errors.InternalWrapf(err, "Failed to call Facebook API /me")
	}
	user := struct {
		ID    string
		Name  string
		Email string
	}{}
	if err := result.Decode(&user); err != nil {
		return "", "", "", errors.InternalWrapf(err, "Failed to decode result of /me")
	}

	return user.ID, user.Name, user.Email, nil
}

func getGoogleOAuthConfig(r *http.Request) oauth2.Config {
	c := googleOAuthConfig
	c.RedirectURL = fmt.Sprintf("%s://%s/oauth/google/callback", config.WebURLScheme(r), r.Host)
	return c
}

func getFacebookOAuthConfig(r *http.Request) oauth2.Config {
	c := facebookOAuthConfig
	c.RedirectURL = fmt.Sprintf("%s://%s/oauth/facebook/callback", config.WebURLScheme(r), r.Host)
	return c
}
