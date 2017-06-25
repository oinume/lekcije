package notifier

import (
	"testing"

	"net/http"

	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
)

func TestSendNotification(t *testing.T) {
	a := assert.New(t)

	h := model.NewTestHelper(t)
	db := h.DB()
	defer db.Close()
	a.Nil(db.DB().Ping())

	mCountries := h.LoadMCountries()
	fetcher := fetcher.NewTeacherLessonFetcher(http.DefaultClient, 1, false, mCountries, nil)
	n := NewNotifier(db, fetcher, true, false)

	a.NotNil(n)
}
