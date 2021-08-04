package flash_message

import (
	"encoding/json"
	"fmt"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/randoms"
)

var _ = fmt.Print

type Kind int

const (
	KindSuccess = iota + 1
	KindInfo
	KindWarning
	KindError
)

func (k Kind) String() string {
	switch k {
	case KindSuccess:
		return "success"
	case KindInfo:
		return "info"
	case KindWarning:
		return "warning"
	case KindError:
		return "error"
	default:
		return ""
	}
}

func (k Kind) ViewStyle() string {
	switch k {
	case KindSuccess:
		return "success"
	case KindInfo:
		return "info"
	case KindWarning:
		return "warning"
	case KindError:
		return "danger"
	default:
		return ""
	}
}

type FlashMessage struct {
	Kind     Kind     `json:"kind"`
	Key      string   `json:"key"`
	Messages []string `json:"messages"`
}

func New(kind Kind, messages ...string) *FlashMessage {
	return &FlashMessage{
		Kind:     kind,
		Key:      randoms.MustNewString(32),
		Messages: messages,
	}
}

func (f *FlashMessage) HasMultipleMessage() bool {
	return len(f.Messages) > 1
}

func (f *FlashMessage) Set() (string, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to json.Marshal()"),
		)
	}
	return string(bytes), nil
}

func (f *FlashMessage) AsURLQueryString() string {
	return "flashMessageKey=" + f.Key
}
