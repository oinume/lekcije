package email

import (
	"os"
	"fmt"
	"strings"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
)

type Sender interface {
	Send(email *Email) error
}

type SendGridSender struct{}

func NewSendGridSender() Sender {
	return &SendGridSender{}
}

func (s *SendGridSender) Send(email *Email) error {
	// TODO: Set HTTP Client
	from := mail.NewEmail(email.From.Name, email.From.Address)
	to := mail.NewEmail(email.To.Name, email.To.Address)
	content := mail.NewContent("text/html", strings.Replace(email.Body, "\n", "<br>", -1))
	m := mail.NewV3MailInit(from, email.Subject, to, content)

	req := sendgrid.GetRequest(
		os.Getenv("SENDGRID_API_KEY"),
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	resp, err := sendgrid.API(req)
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
