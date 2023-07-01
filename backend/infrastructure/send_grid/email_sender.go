package send_grid

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/domain/model/email"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
)

const (
	sendGridAPIHost = "https://api.sendgrid.com"
	sendGridAPIPath = "/v3/mail/send"
)

var (
	redirectErrorFunc = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	defaultHTTPClient = &http.Client{
		Timeout:       10 * time.Second,
		CheckRedirect: redirectErrorFunc,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 1200 * time.Second,
			}).DialContext,
			IdleConnTimeout:     1200 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				ClientSessionCache: tls.NewLRUClientSessionCache(100),
			},
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)

type sendGridEmailSender struct {
	client    *rest.Client
	appLogger *zap.Logger
}

func NewSendGridEmailSender(httpClient *http.Client, appLogger *zap.Logger) repository.EmailSender {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	client := &rest.Client{
		HTTPClient: httpClient,
	}
	return &sendGridEmailSender{
		client:    client,
		appLogger: appLogger,
	}
}

func (s *sendGridEmailSender) Send(ctx context.Context, email *email.Email) error {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "sendGridEmailSender.Send")
	defer span.End()

	from := mail.NewEmail(email.From.Name, email.From.Address)
	content := mail.NewContent("text/html", strings.Replace(email.BodyString(), "\n", "<br>", -1))
	tos := make([]*mail.Email, len(email.Tos))
	for i, to := range email.Tos {
		tos[i] = mail.NewEmail(to.Name, to.Address)
	}
	m := mail.NewV3MailInit(from, email.Subject, tos[0], content)
	m.Personalizations[0].AddTos(tos[1:]...)
	for k, v := range email.CustomArgs() {
		m.SetCustomArg(k, v)
	}

	req := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), sendGridAPIPath, sendGridAPIHost)
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	//fmt.Printf("--- request ---\n%s\n", string(req.Body))
	resp, err := s.client.Send(req)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to send email by SendGrid"),
		)
	}
	//fmt.Printf("--- response ---\nstatus=%d\n%s\n", resp.StatusCode, resp.Body)
	// No need to resp.Body.Close(). It's a string
	if resp.StatusCode >= 300 {
		message := fmt.Sprintf(
			"Failed to send email by SendGrid: statusCode=%v, body=%v",
			resp.StatusCode, strings.Replace(resp.Body, "\n", "\\n", -1),
		)
		s.appLogger.Error(message) // TODO: remove and log in caller
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to send email by SendGrid: statusCode=%v", resp.StatusCode),
		)
	}

	return nil
}
