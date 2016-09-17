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
	lessonService           *LessonService
	userService             *UserService
	userGoogleService       *UserGoogleService
	userApiTokenService     *UserApiTokenService
)

func TestMain(m *testing.M) {
	dbDsn := os.Getenv("DB_DSN")
	if strings.HasSuffix(dbDsn, "/lekcije") {
		dbDsn = strings.Replace(dbDsn, "/lekcije", "/lekcije_test", 1)
	}
	var err error
	db, err = OpenDB(dbDsn)
	if err != nil {
		panic(err)
	}

	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
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
