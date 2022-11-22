package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/usecase"
)

type BatchServer struct{}

func NewBatchServer(
	appLogger *zap.Logger,
	db *gorm.DB,
	errorRecorder *usecase.ErrorRecorder,
	fetcher repository.LessonFetcher,
	lessonUsecase *usecase.Lesson,
	sender emailer.Sender,
) *BatchServer {
	return &BatchServer{}
}

/*
	concurrency          = flagSet.Int("concurrency", 1, "Concurrency of fetcher")
	dryRun               = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
	fetcherCache         = flagSet.Bool("fetcher-cache", false, "Cache teacher and lesson data in Fetcher")
	notificationInterval = flagSet.Int("notification-interval", 0, "Notification interval")
	sendEmail            = flagSet.Bool("send-email", true, "Flag to send email")
	logLevel             = flagSet.String("log-level", "info", "Log level")
*/

func (s *BatchServer) notifier(w http.ResponseWriter, r *http.Request) {
	if err := s.validateRequest(r); err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	type Params struct {
		Concurrency int `json:"concurrency"`
	}

	params := &Params{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}
	notifier := usecase.NewNotifier()
}

func (s *BatchServer) validateRequest(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("invalid content type")
	}
}
