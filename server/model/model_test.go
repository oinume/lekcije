package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/util"
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
	teacherService          *TeacherService
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
	teacherService = NewTeacherService(db)
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

type TestHelper struct {
	t  *testing.T
	db *gorm.DB
}

func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t: t,
	}
}

func (h *TestHelper) DB() *gorm.DB {
	if h.db != nil {
		return h.db
	}
	bootstrap.CheckCLIEnvVars()
	testDBURL = ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL)
	db, err := OpenDB(testDBURL, 1, true) // TODO: env
	if err != nil {
		h.t.Fatalf("Failed to OpenDB(): err=%v", err)
	}
	h.db = db
	return db
}

func (h *TestHelper) CreateUser(name, email string) *User {
	db := h.DB()
	user, err := NewUserService(db).Create(name, email)
	if err != nil {
		h.t.Fatalf("Failed to CreateUser(): err=%v", err)
	}
	return user
}

func (h *TestHelper) CreateRandomUser() *User {
	name := util.RandomString(16)
	return h.CreateUser(name, name+"@example.com")
}

func (h *TestHelper) CreateTeacher(id uint32, name string) *Teacher {
	db := h.DB()
	teacher := &Teacher{
		ID:     id,
		Name:   name,
		Gender: "female",
	}
	if err := NewTeacherService(db).CreateOrUpdate(teacher); err != nil {
		h.t.Fatalf("Failed to CreateTeacher(): err=%v", err)
	}
	return teacher
}
