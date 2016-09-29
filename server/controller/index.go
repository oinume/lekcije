package controller

import (
	"fmt"
	"net/http"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/controller/flash_message"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

var _ = fmt.Print

type navigationItem struct {
	Text string
	URL  string
}

type commonTemplateData struct {
	StaticURL         string
	GoogleAnalyticsID string
	CurrentURL        string
	NavigationItems   []navigationItem
}

var loggedInNavigationItems = []navigationItem{
	{"Home", "/"},
	{"Setting", "/me/setting"},
	{"Logout", "/logout"},
}

var loggedOutNavigationItems = []navigationItem{
	{"Home", "/"},
}

func getCommonTemplateData(currentURL string, loggedIn bool) commonTemplateData {
	data := commonTemplateData{
		StaticURL:         config.StaticURL(),
		GoogleAnalyticsID: config.GoogleAnalyticsID(),
		CurrentURL:        currentURL,
	}
	if loggedIn {
		data.NavigationItems = loggedInNavigationItems
	} else {
		data.NavigationItems = loggedOutNavigationItems
	}
	return data
}

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := model.GetLoggedInUser(ctx)
	if err == nil {
		indexLogin(w, r, user)
	} else {
		indexLogout(w, r)
	}
}

func indexLogin(w http.ResponseWriter, r *http.Request, user *model.User) {
	ctx := r.Context()
	t := ParseHTMLTemplates(TemplatePath("indexLogin.html"))
	type Data struct {
		commonTemplateData
		Teachers     []*model.Teacher
		FlashMessage *flash_message.FlashMessage
	}
	data := &Data{commonTemplateData: getCommonTemplateData(r.RequestURI, true)}

	flashMessageKey := r.FormValue("flashMessageKey")
	if flashMessageKey != "" {
		flashMessage, _ := flash_message.MustStore(ctx).Load(flashMessageKey)
		data.FlashMessage = flashMessage
	}

	followingTeacherService := model.NewFollowingTeacherService(model.MustDB(ctx))
	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
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

func indexLogout(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("indexLogout.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{commonTemplateData: getCommonTemplateData(r.RequestURI, false)}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	userApiTokenService := model.NewUserApiTokenService(model.MustDB(ctx))
	if err := userApiTokenService.DeleteByUserIDAndToken(user.ID, token); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
