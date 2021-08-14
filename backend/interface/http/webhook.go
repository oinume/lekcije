package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	model2 "github.com/oinume/lekcije/backend/model2c"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/event_logger"
	"github.com/oinume/lekcije/backend/model"
)

type SendGridEventValues struct {
	Timestamp int64  `json:"timestamp"`
	Event     string `json:"event"`
	Email     string `json:"email"`
	SGEventID string `json:"sg_event_id"`
	UserAgent string `json:"useragent"`
	URL       string `json:"url"` // Only when event=click
	// Custom args
	EmailType  string `json:"email_type"`
	UserID     string `json:"user_id"`
	TeacherIDs string `json:"teacher_ids"`
}

func (v *SendGridEventValues) GetUserID() uint32 {
	if id, err := strconv.ParseUint(v.UserID, 10, 32); err == nil {
		return uint32(id)
	}
	return 0
}

func (v *SendGridEventValues) IsEventClick() bool {
	return v.Event == "click"
}

func (v *SendGridEventValues) IsEventOpen() bool {
	return v.Event == "open"
}

func (v *SendGridEventValues) LogToFile(logger *zap.Logger) {
	fields := []zapcore.Field{
		zap.Time("timestamp", time.Unix(v.Timestamp, 0)),
		zap.String("sgEventID", v.SGEventID),
		zap.String("email", v.Email),
	}

	var userID uint32
	if id, err := strconv.ParseUint(v.UserID, 10, 32); err == nil {
		userID = uint32(id)
	}
	if v.EmailType != "" {
		fields = append(fields, zap.String("emailType", v.EmailType))
	}
	if v.TeacherIDs != "" {
		fields = append(fields, zap.String("teacherIDs", v.TeacherIDs))
	}
	if v.IsEventOpen() || v.IsEventClick() {
		fields = append(fields, zap.String("userAgent", v.UserAgent))
	}
	if v.IsEventClick() {
		fields = append(fields, zap.String("url", v.URL))
	}

	event_logger.New(logger).Log(
		userID,
		model2.GAMeasurementEventCategoryEmail,
		v.Event,
		fields...,
	)
}

func (v *SendGridEventValues) LogToDB(db *gorm.DB) error {
	eventLogEmail := &model.EventLogEmail{
		Datetime:   time.Unix(v.Timestamp, 0),
		Event:      v.Event,
		EmailType:  v.EmailType,
		UserID:     v.GetUserID(),
		UserAgent:  v.UserAgent,
		TeacherIDs: v.TeacherIDs,
		URL:        v.URL,
	}
	if v.EmailType == "" {
		eventLogEmail.EmailType = model.EmailTypeNewLessonNotifier
	}
	return model.NewEventLogEmailService(db).Create(eventLogEmail)
}

func (s *server) postAPISendGridEventWebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.postAPISendGridEventWebhook(w, r)
	}
}

func (s *server) postAPISendGridEventWebhook(w http.ResponseWriter, r *http.Request) {
	values := make([]SendGridEventValues, 0, 1000)
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, 0)
		return
	}
	defer r.Body.Close()
	// datetime, user_id, event(enum), event_id(varchar), text

	userService := model.NewUserService(s.db)
	for _, v := range values {
		v.LogToFile(s.accessLogger)
		if err := v.LogToDB(s.db); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, 0)
			return
		}
		if v.EmailType == model.EmailTypeNewLessonNotifier && v.IsEventOpen() {
			if err := userService.UpdateOpenNotificationAt(v.GetUserID(), time.Unix(v.Timestamp, 0).UTC()); err != nil {
				internalServerError(r.Context(), s.errorRecorder, w, err, 0)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "OK")
}
