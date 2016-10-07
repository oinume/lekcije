package model

import (
	"database/sql"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
	"gopkg.in/redis.v4"
)

const (
	dbDatetimeFormat = "2006-01-02 15:04:05"
	maxConnections   = 7
)

type contextKeyDB struct{}
type contextKeyRedis struct{}

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := sql.Open(
		"mysql",
		dsn+"?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo",
	)
	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxConnections)
	db.SetConnMaxLifetime(10 * time.Minute)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to sql.Open()")
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, errors.InternalWrapf(err, "Failed to gorm.Open()")
	}
	gormDB.LogMode(bootstrap.HTTPServerEnvVars.NodeEnv != "production")

	return gormDB, nil
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

type GORMTransactional func(tx *gorm.DB) error

func GORMTransaction(db *gorm.DB, name string, callback GORMTransactional) error {
	tx := db.Begin()
	if tx.Error != nil {
		return errors.InternalWrapf(tx.Error, "Failed to begin transaction: name=%v", name)
	}

	var err error
	success := false
	defer func() {
		if success {
			return
		}
		if err2 := tx.Rollback().Error; err2 != nil {
			err = errors.InternalWrapf(err2, "Failed to rollback transaction: name=%v", name)
		}
	}()

	if err2 := callback(tx); err2 != nil {
		return err2
	}
	if tx.Error != nil {
		return tx.Error
	}
	if err2 := tx.Commit().Error; err2 != nil {
		return errors.InternalWrapf(err2, "Failed to commit transaction: name=%v", name)
	}
	success = true
	return nil
}

/*
func SimpleGormTransaction(db *gorm.DB, name string, fn func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return &TxErrBegin{err: err}
	}

	succeed := false
	defer func() {
		if succeed {
			return
		}

		if rerr := tx.Rollback().Error; rerr != nil {
			err = &TxErrRollback{err: err, cause: err}
		}
	}()

	err = fn(tx)
	if err != nil {
		return
	}

	err = tx.Error
	if err != nil {
		return
	}

	err = tx.Commit().Error
	if err != nil {
		err = &TxErrCommit{err: err}
		return
	}

	succeed = true
	return
}
*/
