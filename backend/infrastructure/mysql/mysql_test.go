package mysql_test

import (
	"os"
	"testing"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/model"
)

var helper *model.TestHelper

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	helper = model.NewTestHelper()
	helper.TruncateAllTables(nil)
	defer func() { _ = helper.DB(nil).Close() }()
	os.Exit(m.Run())
}
