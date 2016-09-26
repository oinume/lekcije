package model

import (
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
	"gopkg.in/redis.v4"
)

const (
	dbDatetimeFormat = "2006-01-02 15:04:05"
)

type contextKeyDB struct{}
type contextKeyRedis struct{}

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(
		"mysql",
		dsn+"?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo",
	)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to gorm.Open()")
	}
	if os.Getenv("NODE_ENV") == "production" {
		db.LogMode(false)
	} else {
		db.LogMode(true)
	}
	return db, nil
}

func OpenDBAndSetToContext(ctx context.Context, dbURL string) (*gorm.DB, context.Context, error) {
	db, err := OpenDB(dbURL)
	if err != nil {
		return nil, nil, err
	}
	c := context.WithValue(ctx, contextKeyDB{}, db)
	return db, c, nil
}

func MustDB(ctx context.Context) *gorm.DB {
	value := ctx.Value(contextKeyDB{})
	if db, ok := value.(*gorm.DB); ok {
		return db
	} else {
		panic("Failed to get *gorm.DB from context")
	}
}

func OpenRedis(redisURL string) (*redis.Client, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}
	password := ""
	if u.User != nil {
		password, _ = u.User.Password()
	}
	client := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: password,
		DB:       0,
	})
	return client, nil
}

func OpenRedisAndSetToContext(ctx context.Context, redisURL string) (*redis.Client, context.Context, error) {
	r, err := OpenRedis(redisURL)
	if err != nil {
		return nil, nil, err
	}
	c := context.WithValue(ctx, contextKeyRedis{}, r)
	return r, c, nil
}

func MustRedis(ctx context.Context) *redis.Client {
	value := ctx.Value(contextKeyRedis{})
	if r, ok := value.(*redis.Client); ok {
		return r
	} else {
		panic("Failed to get *redis.Client from context")
	}
}
