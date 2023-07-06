package gcp

import (
	"encoding/base64"

	"github.com/morikuni/failure"
	"google.golang.org/api/option"
)

type Cleaner func()

func WithCredentialsJSONFromBase64String(value string) (option.ClientOption, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to decode base64 string"))
	}
	return option.WithCredentialsJSON(b), nil
}
