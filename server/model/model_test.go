package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/oinume/lekcije/server/config"
	"github.com/stretchr/testify/require"
)

var (
	_                           = fmt.Print
	helper                      = NewTestHelper()
	testDBURL                   string
	followingTeacherService     *FollowingTeacherService
	lessonService               *LessonService
	lessonStatusLogService      *LessonStatusLogService
	mCountryService             *MCountryService
	mPlanService                *MPlanService
	notificationTimeSpanService *NotificationTimeSpanService
	teacherService              *TeacherService
	userService                 *UserService
	userGoogleService           *UserGoogleService
	userAPITokenService         *UserAPITokenService
	mCountries                  *MCountries
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	db := helper.DB()
	defer db.Close()

	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	lessonStatusLogService = NewLessonStatusLogService(db)
	mCountryService = NewMCountryService(db)
	mPlanService = NewMPlanService(db)
	notificationTimeSpanService = NewNotificationTimeSpanService(db)
	teacherService = NewTeacherService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
	userAPITokenService = NewUserAPITokenService(db)
	mCountries = helper.LoadMCountries()

	helper.TruncateAllTables(db)
	os.Exit(m.Run())
}

func TestOpenRedis(t *testing.T) {
	r := require.New(t)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		r.Fail("Env 'REDIS_URL' required.")
	}
	client, err := OpenRedis(redisURL)
	r.NoError(err)
	defer client.Close()
	r.NoError(client.Ping().Err())
}
