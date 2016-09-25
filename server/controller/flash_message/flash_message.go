package flash_message

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
	"golang.org/x/net/context"
	"gopkg.in/redis.v4"
)

type Kind int

func (k Kind) String() string {
	switch k {
	case KindInfo:
		return "info"
	case KindError:
		return "error"
	default:
		return ""
	}
}

const (
	KindInfo = iota + 1
	KindError
)

type contextKey struct{}

type FlashMessage struct {
	Kind    Kind   `json:"kind"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

func New(kind Kind, message string) *FlashMessage {
	key := fmt.Sprintf("flashMessage-%s", util.RandomString(32))
	return &FlashMessage{
		Kind:    kind,
		Key:     key,
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

type Store interface {
	Load(key string) (*FlashMessage, error)
	Save(key string, value *FlashMessage) error
}

type StoreRedis struct {
	client *redis.Client
}

func NewStoreRedis(client *redis.Client) *StoreRedis {
	return &StoreRedis{client: client}
}

func NewStoreRedisAndSetToContext(
	ctx context.Context, client *redis.Client,
) (Store, context.Context) {
	store := NewStoreRedis(client)
	c := context.WithValue(ctx, contextKey{}, store)
	return store, c
}

func MustStore(ctx context.Context) Store {
	value := ctx.Value(contextKey{})
	if store, ok := value.(Store); ok {
		return store
	} else {
		panic("Failed to get Store from context")
	}
}

func (s *StoreRedis) Load(key string) (*FlashMessage, error) {
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

func (s *StoreRedis) Save(key string, value *FlashMessage) error {
	return s.SaveWithExpiration(key, value, time.Hour*24)
}

func (s *StoreRedis) SaveWithExpiration(key string, value *FlashMessage, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(key, string(bytes), expiration).Err()
}
