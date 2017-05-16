package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/oinume/lekcije/server/logger"
	"github.com/uber-go/zap"
)

/*
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

*/
func PostAPISendGridEventWebhook(w http.ResponseWriter, r *http.Request) {
	type EventParam struct {
		Timestamp int64  `json:"timestamp"`
		Event     string `json:"event"`
		Email     string `json:"email"`
		SGEventID string `json:"sg_event_id"`
		UserAgent string `json:"useragent"`
		URL       string `json:"url"` // Only when event=click
		// Custom args
		UserID     string `json:"user_id"`
		TeacherIDs string `json:"teacher_ids"`
	}
	params := make([]EventParam, 0, 1000)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		InternalServerError(w, err)
		return
	}
	defer r.Body.Close()
	// datetime, user_id, event(enum), event_id(varchar), text

	for _, p := range params {
		fields := []zap.Field{
			zap.Time("timestamp", time.Unix(p.Timestamp, 0)),
			zap.String("event", p.Event),
			zap.String("sgEventID", p.SGEventID),
			zap.String("email", p.Email),
		}

		if userID, err := strconv.ParseUint(p.UserID, 10, 32); err == nil {
			fields = append(fields, zap.Uint("userID", uint(userID)))
		}
		if p.TeacherIDs != "" {
			fields = append(fields, zap.String("teacherIDs", p.TeacherIDs))
		}

		if p.Event == "open" || p.Event == "click" {
			fields = append(fields, zap.String("userAgent", p.UserAgent))
		}
		if p.Event == "click" {
			fields = append(fields, zap.String("url", p.URL))
		}

		logger.Access.Info("event", fields...)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
