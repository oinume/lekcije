package model

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/goenum"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
)

const (
	contextKeyDb           = "db"
	contextKeyLoggedInUser = "loggedInUser"
)

type User struct {
	Id        uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*User) TableName() string {
	return "user"
}

func (*UserApiToken) TableName() string {
	return "user_api_token"
}

type UserApiToken struct {
	Token     string `gorm:"primary_key;AUTO_INCREMENT"`
	UserId    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthGoogle struct {
	UserId      uint32    `db:"user_id",gorm:"primary_key"`
	AccessToken string    `db:"access_token"`
	IdToken     string    `db:"id_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (*AuthGoogle) TableName() string {
	return "auth_google"
}

type Teacher struct {
	Id        uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Teacher) TableName() string {
	return "teacher"
}

const teacherUrlBase = "http://eikaiwa.dmm.com/teacher/index/%v/"

var (
	idRegexp         = regexp.MustCompile(`^[\d]+$`)
	teacherUrlRegexp = regexp.MustCompile(`https?://eikaiwa.dmm.com/teacher/index/([\d]+)`)
)

func NewTeacher(id uint32) *Teacher {
	return &Teacher{Id: id}
}

func NewTeacherFromIdOrUrl(idOrUrl string) (*Teacher, error) {
	if idRegexp.MatchString(idOrUrl) {
		id, _ := strconv.ParseUint(idOrUrl, 10, 32)
		return NewTeacher(uint32(id)), nil
	} else if group := teacherUrlRegexp.FindStringSubmatch(idOrUrl); len(group) > 0 {
		id, _ := strconv.ParseUint(group[1], 10, 32)
		return NewTeacher(uint32(id)), nil
	}
	return nil, fmt.Errorf("Failed to parse idOrUrl: %s", idOrUrl)
}

func (t *Teacher) Url() string {
	return fmt.Sprintf(teacherUrlBase, t.Id)
}

type Lesson struct {
	TeacherId uint32
	Datetime  time.Time
	Status    string // TODO: enum
}

func (*Lesson) TableName() string {
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

type FollowingTeacher struct {
	UserId    uint32 `gorm:"primary_key"`
	TeacherId uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*FollowingTeacher) TableName() string {
	return "following_teacher"
}

func Open() (*gorm.DB, error) {
	dbDsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%v?charset=utf8mb4&parseTime=true&loc=UTC", dbDsn),
	)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to gorm.Open()")
	}
	return db, nil
}

func OpenAndSetToContext(ctx context.Context) (*gorm.DB, context.Context, error) {
	db, err := Open()
	db.LogMode(true)
	if err != nil {
		return nil, nil, err
	}
	c := context.WithValue(ctx, contextKeyDb, db)
	return db, c, nil
}

func MustDb(ctx context.Context) *gorm.DB {
	value := ctx.Value(contextKeyDb)
	if db, ok := value.(*gorm.DB); ok {
		return db
	} else {
		panic("Failed to get *gorm.DB from context")
	}
}

func FindLoggedInUserAndSetToContext(token string, ctx context.Context) (*User, context.Context, error) {
	db := MustDb(ctx)
	user := &User{}
	sql := `
		SELECT * FROM user AS u
		INNER JOIN user_api_token AS uat ON u.id = uat.user_id
		WHERE uat.token = ?
		`
	result := db.Model(&User{}).Raw(strings.TrimSpace(sql), token).Scan(user)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil, errors.NotFoundWrapf(result.Error, "Failed to find user: token=%s", token)
		}
		return nil, nil, errors.InternalWrapf(result.Error, "find user: token=%s", token)
	}
	c := context.WithValue(ctx, contextKeyLoggedInUser, user)
	return user, c, nil
}

func GetLoggedInUser(ctx context.Context) (*User, error) {
	value := ctx.Value(contextKeyLoggedInUser)
	if user, ok := value.(*User); ok {
		return user, nil
	}
	return nil, errors.NotFoundf("Logged in user not found in context")
}

func MustLoggedInUser(ctx context.Context) *User {
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}
