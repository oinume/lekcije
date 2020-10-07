package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oinume/lekcije/server/ga_measurement"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	google_auth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/registration_email"
	"github.com/oinume/lekcije/server/util"
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
	return fmt.Sprintf("oauthError: unknown error: %d", int(e))
}

func (s *server) oauthGoogleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.oauthGoogle(w, r)
	}
}

func (s *server) oauthGoogle(w http.ResponseWriter, r *http.Request) {
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

func (s *server) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if err := checkState(r); err != nil {
		internalServerError(s.appLogger, w, err, 0)
		return
	}
	token, idToken, err := exchange(r)
	if err != nil {
		if err == oauthErrorAccessDenied {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		internalServerError(s.appLogger, w, err, 0)
		return
	}
	googleID, name, email, err := getGoogleUserInfo(token, idToken)
	if err != nil {
		internalServerError(s.appLogger, w, err, 0)
		return
	}

	userService := model.NewUserService(s.db)
	user, err := userService.FindByGoogleID(googleID)
	userCreated := false
	if err == nil {
		go s.sendGAMeasurementEvent(
			r.Context(),
			ga_measurement.CategoryUser,
			"login",
			fmt.Sprint(user.ID),
			0,
			user.ID,
		)
	} else {
		if !errors.IsNotFound(err) {
			internalServerError(s.appLogger, w, err, 0)
			return
		}
		// Couldn't find user for the googleID, so create a new user
		errTx := model.GORMTransaction(s.db, "OAuthGoogleCallback", func(tx *gorm.DB) error {
			var errCreate error
			user, _, errCreate = userService.CreateWithGoogle(name, email, googleID)
			return errCreate
		})
		if errTx != nil {
			internalServerError(s.appLogger, w, errTx, 0)
			return
		}
		userCreated = true
		go s.sendGAMeasurementEvent(
			r.Context(),
			ga_measurement.CategoryUser,
			"create",
			fmt.Sprint(user.ID),
			0,
			user.ID,
		)
	}

	userAPITokenService := model.NewUserAPITokenService(s.db)
	userAPIToken, err := userAPITokenService.Create(user.ID)
	if err != nil {
		internalServerError(s.appLogger, w, err, user.ID)
		return
	}

	if userCreated {
		// Send registration email
		go func(user *model.User) {
			sender := registration_email.NewEmailSender(s.senderHTTPClient, s.appLogger)
			if err := sender.Send(r.Context(), user); err != nil {
				s.appLogger.Error(
					"Failed to send registration email",
					zap.String("email", user.Email), zap.Error(err),
				)
				util.SendErrorToRollbar(err, fmt.Sprint(user.ID))
			}
		}(user)
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
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage(fmt.Sprintf(
				"Failed to get cookie oauthState: userAgent=%v, remoteAddr=%v",
				r.UserAgent(), getRemoteAddress(r),
			)),
		)
	}
	if state != oauthState.Value {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("state mismatch"),
		)
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
	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to exchange"),
		)
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", errors.NewInternalError(
			errors.WithMessage("Failed to get id_token"),
		)
	}
	return token, idToken, nil
}

// Returns userId, name, email, error
func getGoogleUserInfo(token *oauth2.Token, idToken string) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	oauth2Client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	service, err := google_auth2.NewService(
		context.Background(),
		// TODO: Not sure which is correct
		//option.WithTokenSource(oauth2.StaticTokenSource(token)),
		option.WithHTTPClient(oauth2Client),
	)
	if err != nil {
		return "", "", "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to create oauth2.client"),
		)
	}

	userinfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return "", "", "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to get userinfo"),
		)
	}

	return userinfo.Id, userinfo.Name, userinfo.Email, nil
}

func getGoogleOAuthConfig(r *http.Request) oauth2.Config {
	c := googleOAuthConfig
	host := r.Header.Get("X-Original-Host") // For ngrok
	if host == "" {
		host = r.Host
	}
	c.RedirectURL = fmt.Sprintf("%s://%s/oauth/google/callback", config.DefaultVars.WebURLScheme(r), host)
	return c
}
