package ga_measurement

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	ga "github.com/jpillora/go-ogle-analytics"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/event_logger"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model2"
)

var defaultHTTPClient = &http.Client{
	Transport: &logger.LoggingHTTPTransport{DumpHeaderBody: true},
	Timeout:   time.Second * 7,
}

type Client interface {
	SendEvent(
		ctx context.Context,
		values *model2.GAMeasurementEvent,
		category,
		action,
		label string,
		value int64,
		userID uint32,
	) error
}

type client struct {
	eventLogger *event_logger.Logger
	httpClient  *http.Client
}

var (
	_ Client = (*client)(nil)
	_ Client = (*fakeClient)(nil)
)

func NewClient(
	httpClient *http.Client,
	eventLogger *event_logger.Logger,
) *client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	return &client{
		eventLogger: eventLogger,
		httpClient:  httpClient,
	}
}

func (c *client) SendEvent(
	ctx context.Context,
	values *model2.GAMeasurementEvent,
	category,
	action,
	label string,
	value int64,
	userID uint32,
) error {
	gaClient, err := ga.NewClient(os.Getenv("GOOGLE_ANALYTICS_ID"))
	if err != nil {
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

	//logFields := []zapcore.Field{
	//	zap.String("category", category),
	//	zap.String("action", action),
	//}
	event := ga.NewEvent(category, action)
	if label != "" {
		event.Label(label)
		//logFields = append(logFields, zap.String("label", label))
	}
	if value != 0 {
		event.Value(value)
		//logFields = append(logFields, zap.Int64("value", value))
	}
	if userID != 0 {
		gaClient.UserID(fmt.Sprint(userID))
		//logFields = append(logFields, zap.Uint("userID", uint(userID)))
	}

	c.eventLogger.Log(userID, category, action, zap.String("label", label), zap.Int64("value", value))
	if err := gaClient.Send(event); err != nil {
		return errors.NewInternalError(
			errors.WithMessage("gaClient.Record failed"),
			errors.WithError(err),
		)
	}

	return nil
}

type fakeClient struct{}

func NewFakeClient() *fakeClient {
	return &fakeClient{}
}

func (fc *fakeClient) SendEvent(
	ctx context.Context,
	values *model2.GAMeasurementEvent,
	category,
	action,
	label string,
	value int64,
	userID uint32,
) error {
	return nil
}
