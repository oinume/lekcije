package flash_message

import (
	"testing"

	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	a := assert.New(t)
	client, err := model.OpenRedis("redis://h:@192.168.99.100:16379")
	a.NoError(err)
	store := NewStoreRedis(client)
	err = store.Save("key1", &FlashMessage{Kind: KindInfo, Message: "Your operation succeeded!"})
	a.NoError(err)
}
