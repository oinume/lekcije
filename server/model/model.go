package model

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
)

const (
	contextKeyDb = "db"
)

type AuthGoogle struct {
	UserId      uint32 `gorm:"primary_key"`
	AccessToken string
	IdToken     string // TODO: GoogleID？
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*AuthGoogle) TableName() string {
	return "auth_google"
}

func Open() (*gorm.DB, error) {
	dbDsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(
		"mysql",
		dbDsn+"?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo",
	)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to gorm.Open()")
	}
	db.LogMode(true) // TODO: off in production
	return db, nil
}

func OpenAndSetToContext(ctx context.Context) (*gorm.DB, context.Context, error) {
	db, err := Open()
	if err != nil {
		return nil, nil, err
	}
	c := context.WithValue(ctx, contextKeyDb, db)
	attachDbToRepo(db)
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

func attachDbToRepo(db *gorm.DB) {
	FollowingTeacherRepo.db = db
	UserApiTokenRepo.db = db
}
