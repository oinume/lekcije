package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/util"
)

type TestHelper struct {
	dbURL           string
	db              *gorm.DB
	mCountryService *MCountryService
}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) DB() *gorm.DB {
	if h.db != nil {
		return h.db
	}
	bootstrap.CheckCLIEnvVars()
	h.dbURL = ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL())
	db, err := OpenDB(h.dbURL, 1, false /* TODO: env */)
	if err != nil {
		panic(fmt.Sprintf("Failed to OpenDB(): err=%v", err))
	}
	h.db = db
	return db
}

func (h *TestHelper) getDBName(dbURL string) string {
	if index := strings.LastIndex(dbURL, "/"); index != -1 {
		return dbURL[index+1:]
	}
	return ""
}

func (h *TestHelper) LoadAllTables(db *gorm.DB) []string {
	type Table struct {
		Name string `gorm:"column:table_name"`
	}
	tables := []Table{}
	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = ?"
	if err := db.Raw(sql, h.getDBName(h.dbURL)).Scan(&tables).Error; err != nil {
		panic(fmt.Sprintf("Failed to select table names: err=%v", err))
	}

	tableNames := []string{}
	for _, t := range tables {
		if t.Name == "goose_db_version" {
			continue
		}
		tableNames = append(tableNames, t.Name)
	}
	return tableNames
}

func (h *TestHelper) TruncateAllTables(db *gorm.DB) {
	tables := h.LoadAllTables(db)
	for _, t := range tables {
		if strings.HasPrefix(t, "m_") {
			continue
		}
		if err := db.Exec("TRUNCATE TABLE " + t).Error; err != nil {
			panic(fmt.Sprintf("Failed to truncate table: table=%v, err=%v", t, err))
		}
	}
}

func (h *TestHelper) CreateUser(name, email string) *User {
	db := h.DB()
	user, err := NewUserService(db).Create(name, email)
	if err != nil {
		panic(fmt.Sprintf("Failed to CreateUser(): err=%v", err))
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
		panic(fmt.Sprintf("Failed to CreateUserGoogle(): %v", err))
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
		panic(fmt.Sprintf("Failed to CreateTeacher(): err=%v", err))
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
		panic(fmt.Sprintf("Failed to MCountryService.LoadAll(): err=%v", err))
	}
	return mCountries
}

func (h *TestHelper) CreateFollowingTeacher(userID uint32, teacher *Teacher) *FollowingTeacher {
	now := time.Now()
	ft, err := NewFollowingTeacherService(h.DB()).FollowTeacher(userID, teacher, now)
	if err != nil {
		panic(fmt.Sprintf("Failed to FollowTeacher(): err=%v", err))
	}
	return ft
}
