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
	db                      *gorm.DB
	testDBURL               string
	followingTeacherService *FollowingTeacherService
	lessonService           *LessonService
	mCountryService         *MCountryService
	planService             *PlanService
	userService             *UserService
	userGoogleService       *UserGoogleService
	userAPITokenService     *UserAPITokenService
	mCountries              *MCountries
)

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	testDBURL = ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL)
	var err error
	db, err = OpenDB(testDBURL, 1, true) // TODO: env
	if err != nil {
		panic(err)
	}

	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	mCountryService = NewMCountryService(db)
	planService = NewPlanService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
	userAPITokenService = NewUserAPITokenService(db)

	mCountries, err = mCountryService.LoadAll()
	if err != nil {
		panic(err)
	}

	if err := TruncateAllTables(db, GetDBName(testDBURL)); err != nil {
		panic(err)
	}

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
