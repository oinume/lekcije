package controller

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	helper = model.NewTestHelper()
)

func TestMain(m *testing.M) {
	db := helper.DB()
	defer db.Close()
	helper.TruncateAllTables(db)
	os.Exit(m.Run())
}

func TestPostAPISendGridEventWebhook(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	reqBody := strings.NewReader(`
[
  {
    "email": "oinume@gmail.com",
    "timestamp": 1492528264,
    "teacher_ids": "16944",
    "ip": "10.43.18.4",
    "sg_event_id": "MzJiZWY5YjYtZjQ5Mi00OWM1LTliYWItNzE2ZTZhZDAxYWFm",
    "user_id": "1",
    "sg_message_id": "UjFrt84CTBGkXCjrVwTULw.filter0041p1las1-28480-58F62C73-2B.0",
    "useragent": "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)",
    "event": "open"
  }
]
	`)
	req, err := http.NewRequest("POST", "/api/sendGrid/eventWebhook", reqBody)
	r.Nil(err)
	ctx := context_data.SetDB(req.Context(), helper.DB())
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(PostAPISendGridEventWebhook)
	handler.ServeHTTP(resp, req)

	a.Equal(http.StatusOK, resp.Code)
	// TODO: more assertions
}
