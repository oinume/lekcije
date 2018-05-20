package flash_message

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
	"gopkg.in/redis.v4"
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
		Key:      util.RandomString(32),
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

type Store interface {
	Load(key string) (*FlashMessage, error)
	Save(value *FlashMessage) error
}

type StoreRedis struct {
	client *redis.Client
}

func NewStoreRedis(client *redis.Client) *StoreRedis {
	return &StoreRedis{client: client}
}

func (s *StoreRedis) Load(key string) (*FlashMessage, error) {
	value, err := s.client.Get(s.getKey(key)).Result()
	if err != nil {
		return nil, err
	}
	v := &FlashMessage{}
	if err := json.Unmarshal([]byte(value), v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *StoreRedis) Save(value *FlashMessage) error {
	return s.SaveWithExpiration(value, time.Hour*24)
}

func (s *StoreRedis) SaveWithExpiration(value *FlashMessage, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(s.getKey(value.Key), string(bytes), expiration).Err()
}

// Append prefix to key
func (s *StoreRedis) getKey(key string) string {
	return "flashMessage:" + key
}
