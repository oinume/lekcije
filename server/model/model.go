package model

import (
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
	"gopkg.in/redis.v4"
)

const (
	contextKeyDB     = "db"
	dbDatetimeFormat = "2006-01-02 15:04:05"
)

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(
		"mysql",
		dsn+"?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo",
	)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to gorm.Open()")
	}
	db.LogMode(true) // TODO: off in production
	return db, nil
}

func OpenDBAndSetToContext(ctx context.Context, dsn string) (*gorm.DB, context.Context, error) {
	db, err := OpenDB(dsn)
	if err != nil {
		return nil, nil, err
	}
	c := context.WithValue(ctx, contextKeyDB, db)
	return db, c, nil
}

func MustDB(ctx context.Context) *gorm.DB {
	value := ctx.Value(contextKeyDB)
	if db, ok := value.(*gorm.DB); ok {
		return db
	} else {
		panic("Failed to get *gorm.DB from context")
	}
}

func OpenRedis(dsn string) (*redis.Client, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	password, _ := u.User.Password()
	client := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
