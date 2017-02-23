package flash_message

import (
	"os"
	"testing"

	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v4"
)

var redisClient *redis.Client
var storeRedis *StoreRedis

func TestMain(m *testing.M) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("Env 'REDIS_URL' is required")
	}
	var err error
	redisClient, err = model.OpenRedis(redisURL)
	if err != nil {
		panic(err)
	}
	storeRedis = NewStoreRedis(redisClient)
	os.Exit(m.Run())
}

func TestStoreRedis_LoadSave(t *testing.T) {
	a := assert.New(t)
	flashMessage := New(KindInfo, "Your operation succeeded!")
	err := storeRedis.Save(flashMessage)
	a.Nil(err)

	v, err := storeRedis.Load(flashMessage.Key)
	a.Nil(err)
	a.Equal("Your operation succeeded!", v.Messages[0])
}
