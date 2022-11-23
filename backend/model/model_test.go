package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/oinume/lekcije/backend/domain/config"
)

var (
	_      = fmt.Print
	helper = NewTestHelper()
	//testDBURL                             string
	eventLogEmailService                  *EventLogEmailService
	followingTeacherService               *FollowingTeacherService
	mCountryService                       *MCountryService
	mPlanService                          *MPlanService
	notificationTimeSpanService           *NotificationTimeSpanService
	statDailyNotificationEventService     *StatDailyNotificationEventService
	statDailyUserNotificationEventService *StatDailyUserNotificationEventService
	statNotifierService                   *StatNotifierService
	teacherService                        *TeacherService
	userService                           *UserService
	mCountries                            *MCountries
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()

	eventLogEmailService = NewEventLogEmailService(db)
	followingTeacherService = NewFollowingTeacherService(db)
	mCountryService = NewMCountryService(db)
	mPlanService = NewMPlanService(db)
	notificationTimeSpanService = NewNotificationTimeSpanService(db)
	statDailyNotificationEventService = NewStatDailyNotificationEventService(db)
	statDailyUserNotificationEventService = NewStatDailyUserNotificationEventService(db)
	statNotifierService = NewStatNotifierService(db)
	teacherService = NewTeacherService(db)
	userService = NewUserService(db)
	mCountries = helper.LoadMCountries(nil)

	helper.TruncateAllTables(nil)
	os.Exit(m.Run())
}
