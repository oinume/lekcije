package emailer

import (
	"fmt"
	"os"
	"strings"
	"net/http"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender interface {
	Send(email *Email) error
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

func (s *SendGridSender) Send(email *Email) error {
	// TODO: Set HTTP Client
	from := mail.NewEmail(email.From.Name, email.From.Address)
	content := mail.NewContent("text/html", strings.Replace(email.BodyString(), "\n", "<br>", -1))
	tos := make([]*mail.Email, len(email.Tos))
	for i, to := range email.Tos {
		tos[i] = mail.NewEmail(to.Name, to.Address)
	}
	m := mail.NewV3MailInit(from, email.Subject, tos[0], content)
	m.Personalizations[0].AddTos(tos[1:]...)

	req := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	resp, err := s.client.API(req)
	if err != nil {
		return errors.InternalWrapf(err, "Failed to send email by sendgrid")
	}
	if resp.StatusCode >= 300 {
		message := fmt.Sprintf(
			"Failed to send email by sendgrid: statusCode=%v, body=%v",
			resp.StatusCode, strings.Replace(resp.Body, "\n", "\\n", -1),
		)
		logger.App.Error(message)
		return errors.InternalWrapf(
			err,
			"Failed to send email by sendgrid: statusCode=%v",
			resp.StatusCode,
		)
	}

	return nil
}

//type NullSender struct {}
//
//func (s *NullSender) Send(email *Email) error {
//	return nil
//}
