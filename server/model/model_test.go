package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/stretchr/testify/assert"
)

var (
	_                       = fmt.Print
	helper                  = NewTestHelper()
	db                      *gorm.DB
	testDBURL               string
	followingTeacherService *FollowingTeacherService
	lessonService           *LessonService
	mCountryService         *MCountryService
	mPlanService            *MPlanService
	teacherService          *TeacherService
	userService             *UserService
	userGoogleService       *UserGoogleService
	userAPITokenService     *UserAPITokenService
	mCountries              *MCountries
)

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	helper.dbURL = ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL())
	db := helper.DB()

	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	mCountryService = NewMCountryService(db)
	mPlanService = NewMPlanService(db)
	teacherService = NewTeacherService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
	userAPITokenService = NewUserAPITokenService(db)
	mCountries = helper.LoadMCountries()

	helper.TruncateAllTables(db)
	os.Exit(m.Run())
}

func TestOpenRedis(t *testing.T) {
	a := assert.New(t)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		a.Fail("Env 'REDIS_URL' required.")
	}
	client, err := OpenRedis(redisURL)
	a.Nil(err)
	defer client.Close()
	a.Nil(client.Ping().Err())
}
