package http

import (
	"io"
	"os"
	"testing"

	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	interfaces "github.com/oinume/lekcije/backend/interface"
	"github.com/oinume/lekcije/backend/logger"
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
		GormDB:              helper.DB(t),
		GAMeasurementClient: ga_measurement.NewFakeClient(),
	})
}
