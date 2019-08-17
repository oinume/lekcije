package http

import (
	"os"
	"testing"

	"github.com/oinume/lekcije/server/ga_measurement"
	"github.com/oinume/lekcije/server/interfaces"
)

func TestMain(m *testing.M) {
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()
	helper.TruncateAllTables(nil)
	_ = os.Chdir("../../..")
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) *server {
	return NewServer(&interfaces.ServerArgs{
		DB:                  helper.DB(t),
		GAMeasurementClient: ga_measurement.NewFakeClient(),
	})
}
