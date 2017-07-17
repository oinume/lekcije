package event_logger

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jpillora/go-ogle-analytics"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/logger"
	"github.com/uber-go/zap"
)

const (
	CategoryEmail            = "email"
	CategoryUser             = "user"
	CategoryFollowingTeacher = "followingTeacher"
)

var gaHTTPClient *http.Client = &http.Client{
	Transport: &logger.LoggingHTTPTransport{DumpHeaderBody: true},
	Timeout:   time.Second * 7,
}

func EventLog(userID uint32, category string, action string, fields ...zap.Field) {
	f := make([]zap.Field, 0, len(fields)+1)
	f = append(
		f,
		zap.String("category", category),
		zap.String("action", action),
		zap.Uint("userID", uint(userID)),
	)
	f = append(f, fields...)
	logger.Access.Info("eventLog", f...)
}

func SendGAMeasurementEvent(req *http.Request, category, action, label string, value int64, userID uint32) {
	gaClient, err := ga.NewClient(os.Getenv("GOOGLE_ANALYTICS_ID"))
	if err != nil {
		logger.App.Warn("ga.NewClient() failed", zap.Error(err))
	}
	gaClient.HttpClient = gaHTTPClient
	gaClient.UserAgentOverride(req.UserAgent())

	gaClient.ClientID(context_data.MustTrackingID(req.Context()))
	gaClient.DocumentHostName(req.Host)
	gaClient.DocumentPath(req.URL.Path)
	gaClient.DocumentTitle(req.URL.Path)
	gaClient.DocumentReferrer(req.Referer())
	gaClient.IPOverride(getRemoteAddress(req))

	logFields := []zap.Field{
		zap.String("category", category),
		zap.String("action", action),
	}
	event := ga.NewEvent(category, action)
	if label != "" {
		event.Label(label)
		logFields = append(logFields, zap.String("label", label))
	}
	if value != 0 {
		event.Value(value)
		logFields = append(logFields, zap.Int64("value", value))
	}
	if userID != 0 {
		gaClient.UserID(fmt.Sprint(userID))
		logFields = append(logFields, zap.Uint("userID", uint(userID)))
	}
	if err := gaClient.Send(event); err == nil {
		logger.App.Debug("SendGAMeasurementEvent() success", logFields...)
		EventLog(userID, category, action, zap.String("label", label), zap.Int64("value", value))
	} else {
		logger.App.Warn("SendGAMeasurementEvent() failed", zap.Error(err))
	}
}

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
