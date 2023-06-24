package gcp

import (
	"encoding/base64"

	"google.golang.org/api/option"

	"github.com/morikuni/failure"
)

func WithCredentialsJSONFromBase64String(value string) (option.ClientOption, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("base64.StdEncoding.DecodeString failed"))
	}
	return option.WithCredentialsJSON(b), nil
}
