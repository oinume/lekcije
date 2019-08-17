package http

import (
	"io"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/ga_measurement"
	"github.com/oinume/lekcije/server/interfaces"
	"github.com/oinume/lekcije/server/logger"
)

func TestMain(m *testing.M) {
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()
	helper.TruncateAllTables(nil)
	_ = os.Chdir("../../..")
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, accessLog io.Writer) *server {
	return NewServer(&interfaces.ServerArgs{
		AccessLogger:        logger.NewAccessLogger(accessLog),
		DB:                  helper.DB(t),
		GAMeasurementClient: ga_measurement.NewFakeClient(),
	})
}
