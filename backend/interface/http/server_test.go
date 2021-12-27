package http

import (
	"context"
	"io"
	"os"
	"testing"

	rollbar_go "github.com/rollbar/rollbar-go"
	"go.uber.org/zap/zapcore"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/infrastructure/rollbar"
	interfaces "github.com/oinume/lekcije/backend/interface"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/usecase"
)

func TestMain(m *testing.M) {
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()
	helper.TruncateAllTables(nil)
	_ = os.Chdir("../..")
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, accessLog io.Writer, appLog io.Writer) *server {
	appLogger := logger.NewAppLogger(appLog, zapcore.InfoLevel)
	gormDB := helper.DB(t)
	rollbarClientMock := &rollbar.ClientMock{
		ErrorWithStackSkipWithExtrasAndContextFunc: func(ctx context.Context, level string, err error, skip int, extras map[string]interface{}) {
			// nop
		},
		SetStackTracerFunc: func(stackTracer rollbar_go.StackTracerFunc) {
			// nop
		},
	}
	return NewServer(
		&interfaces.ServerArgs{
			AccessLogger:        logger.NewAccessLogger(accessLog),
			AppLogger:           appLogger,
			GormDB:              gormDB,
			GAMeasurementClient: ga_measurement.NewFakeClient(),
		},
		usecase.NewErrorRecorder(appLogger, rollbar.NewErrorRecorderRepository(rollbarClientMock)),
		usecase.NewUserAPIToken(mysql.NewUserAPITokenRepository(gormDB.DB())),
	)
}
