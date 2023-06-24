package firebase

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewApp(ctx context.Context, projectID string, serviceAccountKey string) (*firebase.App, error) {
	firebaseConfig := &firebase.Config{
		ProjectID: projectID,
	}
	credential, err := WithCredentialsJSONFromBase64String(serviceAccountKey)
	if err != nil {
		return nil, err
	}
	return firebase.NewApp(ctx, firebaseConfig, credential)
}

func WithCredentialsJSONFromBase64String(value string) (option.ClientOption, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return option.WithCredentialsJSON(b), nil
}
