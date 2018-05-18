package http

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db := helper.DB()
	defer db.Close()
	helper.TruncateAllTables(db)
	_ = os.Chdir("../../..")
	os.Exit(m.Run())
}
