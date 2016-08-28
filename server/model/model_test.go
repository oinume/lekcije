package model

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

var _ = fmt.Print
var db *gorm.DB

func TestMain(m *testing.M) {
	dbDsn := os.Getenv("DB_DSN")
	if strings.HasSuffix(dbDsn, "/lekcije") {
		dbDsn = strings.Replace(dbDsn, "/lekcije", "/lekcije_test", 1)
	}
	var err error
	db, err = Open(dbDsn)
	if err != nil {
		panic(err)
	}
	attachDbToRepo(db)

	for _, t := range []string{"user"} {
		if err := db.Exec("TRUNCATE TABLE " + t).Error; err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}
