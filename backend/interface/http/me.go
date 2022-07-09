package http

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
)

var _ = fmt.Printf

func (s *server) getMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getMe(w, r)
	}
}

func (s *server) getMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.GetTracerProvider().Tracer("http")
	ctx, span := tracer.Start(ctx, "getMe")
	defer span.End()

	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/index.html"))
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, true, user.ID),
		Email:              user.Email,
	}
	if err := t.Execute(w, data); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), user.ID)
		return
	}
}

func (s *server) getMeSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/setting.html"))
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, true, user.ID),
		Email:              user.Email,
	}

	if err := t.Execute(w, data); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), user.ID)
		return
	}
}

func (s *server) getMeLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie(APITokenCookieName)
	if err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to get token cookie"),
		), user.ID)
		return
	}
	token := cookie.Value
	cookieToDelete := &http.Cookie{
		Name:     APITokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	}
	http.SetCookie(w, cookieToDelete)
	if err := s.userAPITokenUsecase.DeleteByUserIDAndToken(ctx, uint(user.ID), token); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	go s.sendGAMeasurementEvent(
		r.Context(),
		model2.GAMeasurementEventCategoryUser,
		"logout",
		fmt.Sprint(user.ID),
		0,
		user.ID,
	)
	http.Redirect(w, r, "/", http.StatusFound)
}
