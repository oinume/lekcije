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
	err := storeRedis.Save("key1", &FlashMessage{Kind: KindInfo, Message: "Your operation succeeded!"})
	a.NoError(err)

	v, err := storeRedis.Load("key1")
	a.NoError(err)
	a.Equal("Your operation succeeded!", v.Message)
}
