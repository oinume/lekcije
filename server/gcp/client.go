package gcp

import (
	"fmt"
	"os"

	"github.com/oinume/lekcije/server/util"
	"google.golang.org/api/option"
)

type Cleaner func()

func WithCredentialsFileFromEnv(envName string) (option.ClientOption, Cleaner, error) {
	serviceAccountKey := os.Getenv(envName)
	if serviceAccountKey == "" {
		return nil, nil, fmt.Errorf("env %s not found", envName)
	}
	return WithCredentialsFileFromBase64String(serviceAccountKey)
}

func WithCredentialsFileFromBase64String(value string) (option.ClientOption, Cleaner, error) {
	f, err := util.GenerateTempFileFromBase64String("", "gcp-", value)
	if err != nil {
		return nil, nil, err
	}
	//defer func() { _ = os.Remove(f.Name()) }()
	return option.WithCredentialsFile(f.Name()), func() { _ = os.Remove(f.Name()) }, nil
}
