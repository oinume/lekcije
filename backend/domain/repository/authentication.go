package repository

//go:generate moq -out=authentication.moq.go . Authentication

import (
	"context"
)

type Authentication interface {
	CreateCustomToken(ctx context.Context, userID uint) (string, error)
}
