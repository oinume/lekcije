package model

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/goenum"
)

type User struct {
	Id        uint32    `db:"id",gorm:"primary_key"`
	Name      string    `db:"name",gorm:"column:name"`
	Email     string    `db:"email",gorm:"column:email"`
	CreatedAt time.Time `db:"created_at",gorm:"column:created_at"`
	UpdatedAt time.Time `db:"updated_at",gorm:"column:updated_at"`
}

func (_ *User) TableName() string {
	return "user"
}

type AuthGoogle struct {
	UserId      uint32    `db:"user_id",gorm:"primary_key"`
	AccessToken string    `db:"access_token"`
	IdToken     string    `db:"id_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (_ *AuthGoogle) TableName() string {
	return "auth_google"
}

type Teacher struct {
	Id        uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (_ *Teacher) TableName() string {
	return "teacher"
}

const teacherUrlBase = "http://eikaiwa.dmm.com/teacher/index/%v/"

func NewTeacher(id uint32) *Teacher {
	return &Teacher{Id: id}
}

func (t *Teacher) Url() string {
	return fmt.Sprintf(teacherUrlBase, t.Id)
}

type Lesson struct {
	TeacherId uint32
	Datetime  time.Time
	Status    string // TODO: enum
}

func (_ *Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherId: %v, Datetime: %v, Status: %v",
		l.TeacherId, l.Datetime, l.Status,
	)
}

type LessonStatus struct {
	Finished   int `goenum:"終了"`
	Reserved   int `goenum:"予約済"`
	Reservable int `goenum:"予約可"`
	Cancelled  int `goenum:"休講"`
}

var LessonStatuses = goenum.EnumerateStruct(&LessonStatus{
	Finished:   1,
	Reserved:   2,
	Reservable: 3,
	Cancelled:  4,
})

func Open() (*gorm.DB, error) {
	dbDsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%v?charset=utf8mb4&parseTime=true&loc=UTC", dbDsn),
	)
	return db, err
}
