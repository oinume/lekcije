package http

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/usecase"
	"github.com/oinume/lekcije/backend/util"
)

const (
	APITokenCookieName   = "apiToken"
	TrackingIDCookieName = "trackingId"
)

func TemplateDir() string {
	koDataPath := os.Getenv("KO_DATA_PATH")
	if koDataPath != "" {
		return filepath.Join(koDataPath, "html")
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, file := filepath.Split(wd)
	switch file {
	case "backend":
		return filepath.Join(wd, "..", "frontend", "html")
	default:
		return filepath.Join(wd, "frontend", "html")
	}
}

func TemplatePath(file string) string {
	return path.Join(TemplateDir(), file)
}

func ParseHTMLTemplates(files ...string) *template.Template {
	f := []string{
		TemplatePath("_base.html"),
	}
	f = append(f, files...)
	return template.Must(template.ParseFiles(f...))
}

func internalServerError(ctx context.Context, errorRecorder *usecase.ErrorRecorder, w http.ResponseWriter, err error, userID uint32) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	sUserID := ""
	if userID == 0 {
		sUserID = fmt.Sprint(userID)
	}
	errorRecorder.Record(ctx, err, sUserID)
	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
	if !config.IsProductionEnv() {
		fmt.Fprintf(w, "----- stacktrace -----\n")
		if e, ok := err.(errors.StackTracer); ok {
			for _, f := range e.StackTrace() {
				_, _ = fmt.Fprintf(w, "%+v\n", f)
			}
		}
		if callStack, ok := failure.CallStackOf(err); ok {
			for _, f := range callStack.Frames() {
				_, _ = fmt.Fprintf(w, "[%s] %v:%v\n", f.Func(), f.File(), f.Line())
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as writeJSON" }`, http.StatusInternalServerError)
		return
	}
}

type commonTemplateData struct {
	StaticURL         string
	GoogleAnalyticsID string
	CurrentURL        string
	CanonicalURL      string
	TrackingID        string
	IsUserAgentPC     bool
	IsUserAgentSP     bool
	IsUserAgentTablet bool
	UserID            string
	NavigationItems   []navigationItem
	ServiceEnv        string
}

type navigationItem struct {
	Text      string
	URL       string
	NewWindow bool
}

var loggedInNavigationItems = []navigationItem{
	{"ホーム", "/me", false},
	{"設定", "/me/setting", false},
	{"お問い合わせ", "https://goo.gl/forms/CIGO3kpiQCGjtFD42", true},
	{"ログアウト", "/me/logout", false},
}

var loggedOutNavigationItems = []navigationItem{
	{"ホーム", "/", false},
}

func (s *server) getCommonTemplateData(req *http.Request, loggedIn bool, userID uint32) commonTemplateData {
	canonicalURL := fmt.Sprintf("%s://%s%s", config.WebURLScheme(req), req.Host, req.RequestURI)
	canonicalURL = (strings.SplitN(canonicalURL, "?", 2))[0] // TODO: use url.Parse
	data := commonTemplateData{
		StaticURL:         config.StaticURL(),
		GoogleAnalyticsID: config.DefaultVars.GoogleAnalyticsID,
		CurrentURL:        req.RequestURI,
		CanonicalURL:      canonicalURL,
		IsUserAgentPC:     util.IsUserAgentPC(req),
		IsUserAgentSP:     util.IsUserAgentSP(req),
		IsUserAgentTablet: util.IsUserAgentTablet(req),
		ServiceEnv:        config.DefaultVars.ServiceEnv,
	}

	if loggedIn {
		data.NavigationItems = loggedInNavigationItems
	} else {
		data.NavigationItems = loggedOutNavigationItems
	}
	data.TrackingID = context_data.MustTrackingID(req.Context())
	if userID != 0 {
		data.UserID = fmt.Sprint(userID)
	}

	return data
}

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
