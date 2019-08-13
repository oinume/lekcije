package model_test

import (
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/model"
)

func ReplaceToTestDBURL(dbURL string) string {
	if strings.HasSuffix(dbURL, "/lekcije") {
		return strings.Replace(dbURL, "/lekcije", "/lekcije_test", 1)
	}
	return dbURL
}

type TestHelper struct {
	model.TestHelper
}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) DB(t *testing.T) *gorm.DB {
	return h.TestHelper.DB(t)
}

func (h *TestHelper) GetDBName(dbURL string) string {
	return h.TestHelper.GetDBName(dbURL)
}

func (h *TestHelper) LoadAllTables(t *testing.T, db *gorm.DB) []string {
	return h.TestHelper.LoadAllTables(t, db)
}

func (h *TestHelper) TruncateAllTables(t *testing.T) {
	h.TestHelper.TruncateAllTables(t)
}

func (h *TestHelper) CreateUser(t *testing.T, name, email string) *model.User {
	return h.TestHelper.CreateUser(t, name, email)
}

func (h *TestHelper) CreateRandomUser(t *testing.T) *model.User {
	return h.TestHelper.CreateRandomUser(t)
}

func (h *TestHelper) CreateUserGoogle(t *testing.T, googleID string, userID uint32) *model.UserGoogle {
	return h.TestHelper.CreateUserGoogle(t, googleID, userID)
}

func (h *TestHelper) CreateTeacher(t *testing.T, id uint32, name string) *model.Teacher {
	return h.TestHelper.CreateTeacher(t, id, name)
}

func (h *TestHelper) CreateRandomTeacher(t *testing.T) *model.Teacher {
	return h.TestHelper.CreateRandomTeacher(t)
}

func (h *TestHelper) LoadMCountries(t *testing.T) *model.MCountries {
	return h.TestHelper.LoadMCountries(t)
}

func (h *TestHelper) CreateFollowingTeacher(t *testing.T, userID uint32, teacher *model.Teacher) *model.FollowingTeacher {
	return h.TestHelper.CreateFollowingTeacher(t, userID, teacher)
}
