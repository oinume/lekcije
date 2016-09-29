package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/util"
	"github.com/uber-go/zap"
)

const ApiTokenCookieName = "apiToken"

func TemplateDir() string {
	if util.IsProductionEnv() {
		return "static"
	} else {
		return "src/html"
	}
}

func TemplatePath(file string) string {
	return path.Join(TemplateDir(), file)
}

func ParseHTMLTemplates(files ...string) *template.Template {
	f := []string{
		TemplatePath("_base.html"),
		TemplatePath("_flashMessage.html"),
	}
	f = append(f, files...)
	return template.Must(template.ParseFiles(f...))
}

func InternalServerError(w http.ResponseWriter, err error) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	// TODO: send error to bugsnag or somewhere
	fields := []zap.Field{
		zap.Error(err),
	}
	if e, ok := err.(errors.StackTracer); ok {
		b := &bytes.Buffer{}
		for _, f := range e.StackTrace() {
			fmt.Fprintf(b, "%+v\n", f)
		}
		fields = append(fields, zap.String("stacktrace", b.String()))
	}
	logger.AppLogger.Error("InternalServerError", fields...)

	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
	if !config.IsProductionEnv() {
		fmt.Fprintf(w, "----- stacktrace -----\n")
		if e, ok := err.(errors.StackTracer); ok {
			for _, f := range e.StackTrace() {
				fmt.Fprintf(w, "%+v\n", f)
			}
		}
	}
}

func JSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as JSON" }`, http.StatusInternalServerError)
		return
	}
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
