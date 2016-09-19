package controller

import (
	"encoding/json"

	"github.com/oinume/lekcije/server/errors"
)

type flashMessageKind int

const (
	flashMessageKindInfo = iota
	flashMessageKindError
)

type flashMessage struct {
	Kind    flashMessageKind `json:"kind"`
	Message string           `json:"message"`
}

func newFlashMessage(kind flashMessageKind, message string) *flashMessage {
	return &flashMessage{
		Kind:    kind,
		Message: message,
	}
}

func (f *flashMessage) Set() (string, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed to json.Marshal()")
	}

	return string(bytes), nil
}

type flashMessageStore interface {
	Load(id string) ([]byte, error)
	Save(id string, bytes []byte) error
}

type FlashMessageRedisStore struct{}

func (s *FlashMessageRedisStore) Load() ([]byte, error) {

}

func (s *FlashMessageRedisStore) Save(id string) error {

}
