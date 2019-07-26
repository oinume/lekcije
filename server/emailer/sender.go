package emailer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opencensus.io/trace"
)

const (
	sendGridAPIHost = "https://api.sendgrid.com"
	sendGridAPIPath = "/v3/mail/send"
)

type Sender interface {
	Send(ctx context.Context, email *Email) error
}

type SendGridSender struct {
	client *rest.Client
}

func NewSendGridSender(httpClient *http.Client) Sender {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := &rest.Client{
		HTTPClient: httpClient,
	}
	return &SendGridSender{
		client: client,
	}
}

func (s *SendGridSender) Send(ctx context.Context, email *Email) error {
	_, span := trace.StartSpan(ctx, "SendGridSender.Send")
	defer span.End()

	from := mail.NewEmail(email.From.Name, email.From.Address)
	content := mail.NewContent("text/html", strings.Replace(email.BodyString(), "\n", "<br>", -1))
	tos := make([]*mail.Email, len(email.Tos))
	for i, to := range email.Tos {
		tos[i] = mail.NewEmail(to.Name, to.Address)
	}
	m := mail.NewV3MailInit(from, email.Subject, tos[0], content)
	m.Personalizations[0].AddTos(tos[1:]...)
	for k, v := range email.customArgs {
		m.SetCustomArg(k, v)
	}

	req := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		sendGridAPIPath,
		sendGridAPIHost,
	)
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	//fmt.Printf("--- request ---\n%s\n", string(req.Body))
	resp, err := s.client.API(req)
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
		logger.App.Error(message)
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to send email by SendGrid: statusCode=%v", resp.StatusCode),
		)
	}

	return nil
}

type NoSender struct{}

var _ Sender = (*NoSender)(nil)

func (s *NoSender) Send(ctx context.Context, email *Email) error {
	return nil
}
