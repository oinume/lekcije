package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/domain/model/email"
)

type EmailSender interface {
	Send(ctx context.Context, email *email.Email) error
}

type NopEmailSender struct{}

func NewNopEmailSender() EmailSender {
	return &NopEmailSender{}
}

func (s *NopEmailSender) Send(ctx context.Context, email *email.Email) error {
	return nil
}
