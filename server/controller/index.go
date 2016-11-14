package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

var _ = fmt.Print

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func Index(w http.ResponseWriter, r *http.Request) {
	if _, err := model.GetLoggedInUser(r.Context()); err == nil {
		http.Redirect(w, r, "/me", http.StatusFound)
	} else {
		indexLogout(w, r)
	}
}

func indexLogout(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("index.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, false),
	}

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

	cookie, err := r.Cookie(APITokenCookieName)
	if err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to get token cookie"))
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
	userAPITokenService := model.NewUserAPITokenService(model.MustDB(ctx))
	if err := userAPITokenService.DeleteByUserIDAndToken(user.ID, token); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func RobotsTxt(w http.ResponseWriter, r *http.Request) {
	content := `
User-agent: *
Allow: /
`
	// TODO: sitemap https://www.lekcije.com/sitemap.xml
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, strings.TrimSpace(content))
}

func Terms(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("terms.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, false),
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}
