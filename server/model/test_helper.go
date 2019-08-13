package model

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/util"
)

type TestHelper struct {
	dbURL string
	db    *gorm.DB
	//mCountryService *MCountryService
}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) DB(t *testing.T) *gorm.DB {
	if h.db != nil {
		return h.db
	}
	config.MustProcessDefault()
	h.dbURL = ReplaceToTestDBURL(config.DefaultVars.DBURL())
	db, err := OpenDB(h.dbURL, 1, config.DefaultVars.DebugSQL)
	if err != nil {
		e := fmt.Errorf("OpenDB failed: %v", err)
		if t == nil {
			panic(e)
		} else {
			t.Fatal(e)
		}
	}
	h.db = db
	return db
}

func (h *TestHelper) GetDBName(dbURL string) string {
	if index := strings.LastIndex(dbURL, "/"); index != -1 {
		return dbURL[index+1:]
	}
	return ""
}

func (h *TestHelper) LoadAllTables(t *testing.T, db *gorm.DB) []string {
	type Table struct {
		Name string `gorm:"column:table_name"`
	}
	tables := []Table{}
	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = ?"
	if err := db.Raw(sql, h.GetDBName(h.dbURL)).Scan(&tables).Error; err != nil {
		e := fmt.Errorf("failed to select table names: err=%v", err)
		if t == nil {
			panic(e)
		} else {
			t.Fatal(e)
		}
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

func (h *TestHelper) TruncateAllTables(t *testing.T) {
	//fmt.Printf("TruncateAllTables() called!\n--- stack ---\n%+v\n", errors.NewInternalError().StackTrace())
	db := h.DB(t)
	tables := h.LoadAllTables(t, db)
	for _, table := range tables {
		if strings.HasPrefix(table, "m_") {
			continue
		}
		if err := db.Exec("TRUNCATE TABLE " + table).Error; err != nil {
			e := fmt.Errorf("failed to truncate table: table=%v, err=%v", t, err)
			if t == nil {
				panic(e)
			} else {
				t.Fatal(e)
			}
		}
	}
}

func (h *TestHelper) CreateUser(t *testing.T, name, email string) *User {
	db := h.DB(t)
	user, err := NewUserService(db).Create(name, email)
	if err != nil {
		panic(fmt.Sprintf("Failed to CreateUser(): err=%v", err))
	}
	return user
}

func (h *TestHelper) CreateRandomUser(t *testing.T) *User {
	name := util.RandomString(16)
	return h.CreateUser(t, name, name+"@example.com")
}

func (h *TestHelper) CreateUserGoogle(t *testing.T, googleID string, userID uint32) *UserGoogle {
	userGoogle := &UserGoogle{
		GoogleID: googleID,
		UserID:   userID,
	}
	if err := h.DB(t).Create(userGoogle).Error; err != nil {
		t.Fatal(fmt.Errorf("CreateUserGoogle failed: %v", err))
	}
	return userGoogle
}

func (h *TestHelper) CreateTeacher(t *testing.T, id uint32, name string) *Teacher {
	db := h.DB(t)
	teacher := &Teacher{
		ID:           id,
		Name:         name,
		Gender:       "female",
		Birthday:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		LastLessonAt: time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
	}
	if err := NewTeacherService(db).CreateOrUpdate(teacher); err != nil {
		t.Fatal(fmt.Sprintf("CreateTeacher failed: %v", err))
	}
	return teacher
}

func (h *TestHelper) CreateRandomTeacher(t *testing.T) *Teacher {
	return h.CreateTeacher(t, uint32(util.RandomInt(9999999)), util.RandomString(6))
}

func (h *TestHelper) LoadMCountries(t *testing.T) *MCountries {
	db := h.DB(t)
	// TODO: cache
	mCountries, err := NewMCountryService(db).LoadAll(context.Background())
	if err != nil {
		e := fmt.Errorf("MCountryService.LoadAll failed: %v", err)
		if t == nil {
			panic(e)
		} else {
			t.Fatal(e)
		}
	}
	return mCountries
}

func (h *TestHelper) CreateFollowingTeacher(t *testing.T, userID uint32, teacher *Teacher) *FollowingTeacher {
	now := time.Now()
	ft, err := NewFollowingTeacherService(h.DB(t)).FollowTeacher(userID, teacher, now)
	if err != nil {
		t.Fatal(fmt.Errorf("FollowTeacher failed: %v", err))
	}
	return ft
}
