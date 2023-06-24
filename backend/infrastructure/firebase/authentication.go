package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/domain/repository"
)

type authenticationRepo struct {
	client *auth.Client
}

func NewAuthenticationRepo(client *auth.Client) repository.Authentication {
	return &authenticationRepo{client: client}
}

func (r *authenticationRepo) CreateCustomToken(ctx context.Context, userID uint) (string, error) {
	customToken, err := r.client.CustomTokenWithClaims(ctx, fmt.Sprint(userID), nil)
	if err != nil {
		return "", failure.Wrap(err, failure.Message("failed to create custom token"))
	}
	return customToken, nil
}
