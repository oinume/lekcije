package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oinume/lekcije/server/interfaces"
	"github.com/oinume/lekcije/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	helper = model.NewTestHelper()
)

func TestPostAPISendGridEventWebhook(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser(t)
	const timestamp = 1492528264
	reqBody := strings.NewReader(fmt.Sprintf(`
[
  {
    "email": "oinume@gmail.com",
    "email_type": "new_lesson_notifier",
    "timestamp": %d,
    "teacher_ids": "16944",
    "ip": "10.43.18.4",
    "sg_event_id": "MzJiZWY5YjYtZjQ5Mi00OWM1LTliYWItNzE2ZTZhZDAxYWFm",
    "user_id": "%d",
    "sg_message_id": "UjFrt84CTBGkXCjrVwTULw.filter0041p1las1-28480-58F62C73-2B.0",
    "useragent": "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)",
    "event": "open"
  }
]
	`, timestamp, user.ID))
	req, err := http.NewRequest("POST", "/api/sendGrid/eventWebhook", reqBody)
	r.NoError(err)

	w := httptest.NewRecorder()
	s := NewServer(&interfaces.ServerArgs{
		DB: helper.DB(t),
	})
	handler := s.postAPISendGridEventWebhookHandler()
	handler.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	// Check OpenNotificationAt
	u, err := model.NewUserService(helper.DB(t)).FindByPK(user.ID)
	r.NoError(err)
	r.True(u.OpenNotificationAt.Valid)
	a.Equal(int64(timestamp), u.OpenNotificationAt.Time.Unix())
}
