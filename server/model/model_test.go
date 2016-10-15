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
	userService             *UserService
	userGoogleService       *UserGoogleService
	userApiTokenService     *UserApiTokenService
)

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	testDBURL = ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL)
	var err error
	db, err = OpenDB(testDBURL)
	if err != nil {
		panic(err)
	}

	followingTeacherService = NewFollowingTeacherService(db)
	lessonService = NewLessonService(db)
	userService = NewUserService(db)
	userGoogleService = NewUserGoogleService(db)
	userApiTokenService = NewUserApiTokenService(db)

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
