package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/server/config"
)

var (
	_      = fmt.Print
	helper = NewTestHelper()
	//testDBURL                             string
	eventLogEmailService                  *EventLogEmailService
	followingTeacherService               *FollowingTeacherService
	lessonService                         *LessonService
	lessonStatusLogService                *LessonStatusLogService
	mCountryService                       *MCountryService
	mPlanService                          *MPlanService
	notificationTimeSpanService           *NotificationTimeSpanService
	statDailyNotificationEventService     *StatDailyNotificationEventService
	statDailyUserNotificationEventService *StatDailyUserNotificationEventService
	statNotifierService                   *StatNotifierService
	teacherService                        *TeacherService
	userService                           *UserService
	userGoogleService                     *UserGoogleService
	userAPITokenService                   *UserAPITokenService
	mCountries                            *MCountries
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()

	eventLogEmailService = NewEventLogEmailService(db)
	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	lessonStatusLogService = NewLessonStatusLogService(db)
	mCountryService = NewMCountryService(db)
	mPlanService = NewMPlanService(db)
	notificationTimeSpanService = NewNotificationTimeSpanService(db)
	statDailyNotificationEventService = NewStatDailyNotificationEventService(db)
	statDailyUserNotificationEventService = NewStatDailyUserNotificationEventService(db)
	statNotifierService = NewStatNotifierService(db)
	teacherService = NewTeacherService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
	userAPITokenService = NewUserAPITokenService(db)
	mCountries = helper.LoadMCountries(nil)

	helper.TruncateAllTables(nil)
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
