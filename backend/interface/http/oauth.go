package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"goji.io/v3"
	"goji.io/v3/pat"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	google_auth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/model"
	model2 "github.com/oinume/lekcije/backend/model2c"
	"github.com/oinume/lekcije/backend/registration_email"
	"github.com/oinume/lekcije/backend/usecase"
	"github.com/oinume/lekcije/backend/util"
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

type OAuthServer struct {
	appLogger            *zap.Logger
	gaMeasurementClient  ga_measurement.Client
	gaMeasurementUsecase *usecase.GAMeasurement
	senderHTTPClient     *http.Client
	userUsecase          *usecase.User
	userAPITokenUsecase  *usecase.UserAPIToken
}

func NewOAuthServer(
	appLogger *zap.Logger,
	gaMeasurementClient ga_measurement.Client,
	gaMeasurementUsecase *usecase.GAMeasurement,
	senderHTTPClient *http.Client,
	userUsecase *usecase.User,
	userAPITokenUsecase *usecase.UserAPIToken,
) *OAuthServer {
	return &OAuthServer{
		appLogger:            appLogger,
		gaMeasurementClient:  gaMeasurementClient,
		gaMeasurementUsecase: gaMeasurementUsecase,
		senderHTTPClient:     senderHTTPClient,
		userUsecase:          userUsecase,
		userAPITokenUsecase:  userAPITokenUsecase,
	}
}

func (s *OAuthServer) Setup(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/oauth/google"), s.oauthGoogle)
	mux.HandleFunc(pat.Get("/oauth/google/callback"), s.oauthGoogleCallback)
}

func (s *OAuthServer) oauthGoogle(w http.ResponseWriter, r *http.Request) {
	state := util.RandomString(32)
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
		// TODO: Secure: true
	}
	http.SetCookie(w, cookie)
	c := getGoogleOAuthConfig(r)
	http.Redirect(w, r, c.AuthCodeURL(state), http.StatusFound)
}

func (s *OAuthServer) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()
	user, err := s.userUsecase.FindByGoogleID(ctx, googleID)
	userCreated := false
	if err == nil {
		go func() {
			if err := s.gaMeasurementUsecase.SendEvent(
				r.Context(),
				context_data.MustGAMeasurementEvent(ctx),
				model2.GAMeasurementEventCategoryUser,
				"login",
				fmt.Sprint(user.ID),
				0,
				uint32(user.ID),
			); err != nil {
				// TODO: Define usecase.ErrorRecorder.Warn and Error
				s.appLogger.Warn("SendEvent() failed", zap.Error(err))
			}
		}()
	} else {
		if !errors.IsNotFound(err) {
			internalServerError(s.appLogger, w, err, 0)
			return
		}
		u, _, err := s.userUsecase.CreateWithGoogle(ctx, name, email, googleID)
		if err != nil {
			internalServerError(s.appLogger, w, err, 0)
			return
		}
		userCreated = true
		user = u
		go func() {
			if err := s.gaMeasurementUsecase.SendEvent(
				r.Context(),
				context_data.MustGAMeasurementEvent(ctx),
				model2.GAMeasurementEventCategoryUser,
				"create",
				fmt.Sprint(user.ID),
				0,
				uint32(user.ID),
			); err != nil {
				s.appLogger.Warn("SendEvent() failed", zap.Error(err))
			}
		}()
	}

	userAPIToken, err := s.userAPITokenUsecase.Create(ctx, user.ID)
	if err != nil {
		internalServerError(s.appLogger, w, err, uint32(user.ID))
		return
	}
	s.appLogger.Debug(fmt.Sprintf("userCreated = %v", userCreated))

	if userCreated {
		// TODO: Move to usecase layer
		// Send registration email
		go func() {
			sender := registration_email.NewEmailSender(s.senderHTTPClient, s.appLogger)
			if err := sender.Send(r.Context(), user); err != nil {
				s.appLogger.Error(
					"Failed to send registration email",
					zap.String("email", user.Email), zap.Error(err),
				)
				util.SendErrorToRollbar(err, fmt.Sprint(user.ID))
			}
		}()
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
