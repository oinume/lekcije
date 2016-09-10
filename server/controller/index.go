package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
)

var _ = fmt.Print

type commonTemplateData struct {
	StaticUrl         string
	GoogleAnalyticsId string
}

func getCommonTemplateData() commonTemplateData {
	return commonTemplateData{
		StaticUrl:         config.StaticUrl(),
		GoogleAnalyticsId: config.GoogleAnalyticsId(),
	}
}

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := model.GetLoggedInUser(ctx)
	if err == nil {
		indexLogin(ctx, w, r, user)
	} else {
		indexLogout(ctx, w, r)
	}
}

func indexLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, user *model.User) {
	t := template.Must(template.ParseFiles(
		TemplatePath("_base.html"),
		TemplatePath("indexLogin.html")),
	)
	type Data struct {
		commonTemplateData
		Teachers []*model.Teacher
	}
	data := &Data{commonTemplateData: getCommonTemplateData()}

	teachers, err := model.FollowingTeacherService.FindTeachersByUserId(user.Id)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.Teachers = teachers

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func indexLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		TemplatePath("_base.html"),
		TemplatePath("indexLogout.html")),
	)
	type Data struct {
		commonTemplateData
	}
	data := &Data{commonTemplateData: getCommonTemplateData()}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := model.GetLoggedInUser(ctx)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie(ApiTokenCookieName)
	if err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to get token cookie"))
		return
	}
	token := cookie.Value
	cookieToDelete := &http.Cookie{
		Name:     ApiTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	}
	http.SetCookie(w, cookieToDelete)
	if err := model.UserApiTokenService.DeleteByUserIdAndToken(user.Id, token); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
