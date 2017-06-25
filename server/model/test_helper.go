package model

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/util"
)

type TestHelper struct {
	t               *testing.T
	db              *gorm.DB
	mCountryService *MCountryService
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
	testDBURL := ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL())
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

func (h *TestHelper) CreateUserGoogle(googleID string, userID uint32) *UserGoogle {
	userGoogle := &UserGoogle{
		GoogleID: googleID,
		UserID:   userID,
	}
	if err := h.DB().Create(userGoogle).Error; err != nil {
		h.t.Fatalf("Failed to CreateUserGoogle(): %v", err)
	}
	return userGoogle
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

func (h *TestHelper) CreateRandomTeacher() *Teacher {
	return h.CreateTeacher(uint32(util.RandomInt(99999)), util.RandomString(6))
}

func (h *TestHelper) LoadMCountries() *MCountries {
	db := h.DB()
	// TODO: cache
	mCountries, err := NewMCountryService(db).LoadAll()
	if err != nil {
		h.t.Fatalf("Failed to MCountryService.LoadAll(): err=%v", err)
	}
	return mCountries
}
