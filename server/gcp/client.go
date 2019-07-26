package gcp

import (
	"encoding/base64"

	"github.com/oinume/lekcije/server/errors"
	"google.golang.org/api/option"
)

type Cleaner func()

func WithCredentialsJSONFromBase64String(value string) (option.ClientOption, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, errors.NewInternalError(errors.WithError(err))
	}
	return option.WithCredentialsJSON(b), nil
}
