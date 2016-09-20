package flash_message

import (
	"encoding/json"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"gopkg.in/redis.v4"
)

type Kind int

const (
	KindInfo = iota + 1
	KindError
)

type FlashMessage struct {
	Kind    Kind   `json:"kind"`
	Message string `json:"message"`
}

func New(kind Kind, message string) *FlashMessage {
	return &FlashMessage{
		Kind:    kind,
		Message: message,
	}
}

func (f *FlashMessage) Set() (string, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", errors.InternalWrapf(err, "Failed to json.Marshal()")
	}

	return string(bytes), nil
}

type FlashMessageStore interface {
	Load(key string) ([]byte, error)
	Save(key string, bytes []byte) error
}

type FlashMessageStoreRedis struct {
	client *redis.Client
}

func NewStoreRedis(client *redis.Client) *FlashMessageStoreRedis {
	return &FlashMessageStoreRedis{client: client}
}

func (s *FlashMessageStoreRedis) Load(key string) (*FlashMessage, error) {
	value, err := s.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	v := &FlashMessage{}
	if err := json.Unmarshal([]byte(value), v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *FlashMessageStoreRedis) Save(key string, value *FlashMessage) error {
	return s.SaveWithExpiration(key, value, time.Hour*24)
}

func (s *FlashMessageStoreRedis) SaveWithExpiration(key string, value *FlashMessage, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(key, string(bytes), expiration).Err()
}
