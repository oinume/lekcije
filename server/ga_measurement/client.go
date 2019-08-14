package ga_measurement

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ga "github.com/jpillora/go-ogle-analytics"
	"github.com/oinume/lekcije/server/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultHTTPClient = &http.Client{
	Transport: &logger.LoggingHTTPTransport{DumpHeaderBody: true},
	Timeout:   time.Second * 7,
}

type client struct {
	appLogger   *zap.Logger
	eventLogger *zap.Logger
	httpClient  *http.Client
}

func NewClient(httpClient *http.Client, appLogger *zap.Logger, eventLogger *zap.Logger) *client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	return &client{
		appLogger:   appLogger,
		eventLogger: eventLogger,
		httpClient:  httpClient,
	}
}

func (c *client) SendEvent(
	values *EventValues,
	category,
	action,
	label string,
	value int64,
	userID uint32,
) error {
	gaClient, err := ga.NewClient(os.Getenv("GOOGLE_ANALYTICS_ID"))
	if err != nil {
		c.appLogger.Warn("ga.NewClient() failed", zap.Error(err))
		return err
	}
	gaClient.HttpClient = c.httpClient
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
		// TODO: use event_logger
		c.eventLogger.Log(userID, category, action, zap.String("label", label), zap.Int64("value", value))
	} else {
		logger.App.Warn("SendGAMeasurementEvent() failed", zap.Error(err))
	}
	return nil
}
