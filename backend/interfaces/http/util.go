package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/interfaces/http/flash_message"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/util"
)

const (
	APITokenCookieName   = "apiToken"
	TrackingIDCookieName = "trackingId"
)

func TemplateDir() string {
	return "frontend/html"
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

func internalServerError(appLogger *zap.Logger, w http.ResponseWriter, err error, userID uint32) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	sUserID := ""
	if userID == 0 {
		sUserID = fmt.Sprint(userID)
	}
	util.SendErrorToRollbar(err, sUserID)
	fields := []zapcore.Field{
		zap.Error(err),
	}
	if e, ok := err.(errors.StackTracer); ok {
		b := &bytes.Buffer{}
		for _, f := range e.StackTrace() {
			fmt.Fprintf(b, "%+v\n", f)
		}
		fields = append(fields, zap.String("stacktrace", b.String()))
	}
	if appLogger != nil {
		appLogger.Error("internalServerError", fields...)
	}

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
	FlashMessage      *flash_message.FlashMessage
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
	}

	if loggedIn {
		data.NavigationItems = loggedInNavigationItems
	} else {
		data.NavigationItems = loggedOutNavigationItems
	}
	if flashMessageKey := req.FormValue("flashMessageKey"); flashMessageKey != "" {
		flashMessage, _ := s.flashMessageStore.Load(flashMessageKey)
		data.FlashMessage = flashMessage
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

func authenticateFromContext(ctx context.Context, db *gorm.DB) (*model.User, error) {
	apiToken, err := context_data.GetAPIToken(ctx)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "no api token found")
	}
	userService := model.NewUserService(db)
	user, err := userService.FindLoggedInUser(apiToken)
	if err != nil {
		return nil, twirp.NewError(twirp.Unauthenticated, "no user found")
	}
	return user, nil
}
