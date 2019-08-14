package event_logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	ga "github.com/jpillora/go-ogle-analytics"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	CategoryEmail            = "email"
	CategoryUser             = "user"
	CategoryFollowingTeacher = "followingTeacher"
)

// TODO: Optimize http.Client
var gaHTTPClient = &http.Client{
	Transport: &logger.LoggingHTTPTransport{DumpHeaderBody: true},
	Timeout:   time.Second * 7,
}

// TODO: Rename this package ga_measurement and add EventSender and define methods.

type GAMeasurementEventValues struct {
	UserAgentOverride string
	ClientID          string
	DocumentHostName  string
	DocumentPath      string
	DocumentTitle     string
	DocumentReferrer  string
	IPOverride        string
}

type gaMeasurementEventValuesKey struct{}

func NewGAMeasurementEventValuesFromRequest(req *http.Request) *GAMeasurementEventValues {
	// Ignore if client id is not set
	clientID, _ := context_data.GetTrackingID(req.Context())
	return &GAMeasurementEventValues{
		UserAgentOverride: req.UserAgent(),
		ClientID:          clientID,
		DocumentHostName:  req.Host,
		DocumentPath:      req.URL.Path,
		DocumentTitle:     req.URL.Path,
		DocumentReferrer:  req.Referer(),
		IPOverride:        getRemoteAddress(req),
	}
}

func WithGAMeasurementEventValues(ctx context.Context, v *GAMeasurementEventValues) context.Context {
	return context.WithValue(ctx, gaMeasurementEventValuesKey{}, v)
}

func GetGAMeasurementEventValues(ctx context.Context) (*GAMeasurementEventValues, error) {
	v := ctx.Value(gaMeasurementEventValuesKey{})
	if value, ok := v.(*GAMeasurementEventValues); ok {
		return value, nil
	} else {
		return nil, errors.NewInternalError(
			errors.WithMessage("failed get value from context"),
		)
	}
}

func MustGAMeasurementEventValues(ctx context.Context) *GAMeasurementEventValues {
	v, err := GetGAMeasurementEventValues(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func Log(userID uint32, category string, action string, fields ...zapcore.Field) {
	f := make([]zapcore.Field, 0, len(fields)+1)
	f = append(
		f,
		zap.String("category", category),
		zap.String("action", action),
		zap.Uint("userID", uint(userID)),
	)
	f = append(f, fields...)
	logger.Access.Info("eventLog", f...)
}

// TODO: remove
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

	logFields := []zapcore.Field{
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
		Log(userID, category, action, zap.String("label", label), zap.Int64("value", value))
	} else {
		logger.App.Warn("SendGAMeasurementEvent() failed", zap.Error(err))
	}
}

func SendGAMeasurementEvent2(values *GAMeasurementEventValues, category, action, label string, value int64, userID uint32) {
	gaClient, err := ga.NewClient(os.Getenv("GOOGLE_ANALYTICS_ID"))
	if err != nil {
		logger.App.Warn("ga.NewClient() failed", zap.Error(err))
	}
	gaClient.HttpClient = gaHTTPClient
	gaClient.UserAgentOverride(values.UserAgentOverride)
	gaClient.ClientID(values.ClientID)
	gaClient.DocumentHostName(values.DocumentHostName)
	gaClient.DocumentPath(values.DocumentPath)
	gaClient.DocumentTitle(values.DocumentTitle)
	gaClient.DocumentReferrer(values.DocumentReferrer)
	gaClient.IPOverride(values.IPOverride)

	logFields := []zapcore.Field{
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
		Log(userID, category, action, zap.String("label", label), zap.Int64("value", value))
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
