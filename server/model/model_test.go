package model

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

var (
	_                       = fmt.Print
	db                      *gorm.DB
	followingTeacherService *FollowingTeacherService
	userService             *UserService
	userApiTokenService     *UserApiTokenService
)

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
	attachDbToService(db)
	userService = NewUserService(db)
	userApiTokenService = NewUserApiTokenService(db)

	tables := []string{
		"following_teacher", "lesson",
		"user", "user_api_token", "user_google", "teacher",
	}
	for _, t := range tables {
		if err := db.Exec("TRUNCATE TABLE " + t).Error; err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}
