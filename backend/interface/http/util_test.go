package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/morikuni/failure"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/usecase"
)

func TestInternalServerError(t *testing.T) {
	errorRecorder := usecase.NewErrorRecorder(
		logger.NewAppLogger(new(bytes.Buffer), zapcore.InfoLevel),
		&repository.NopErrorRecorder{},
	)
	w := httptest.NewRecorder()
	err := failure.Wrap(fmt.Errorf("new error"))
	internalServerError(context.Background(), errorRecorder, w, err, 1)

	if body := w.Body.String(); !strings.Contains(body, "code.Internal: new error") {
		t.Fatalf("internalServerError response body is invalid: %v", body)
	}
}
