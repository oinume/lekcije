package http

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()
	helper.TruncateAllTables(nil)
	_ = os.Chdir("../../..")
	os.Exit(m.Run())
}
