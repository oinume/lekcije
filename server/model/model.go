package model

import (
	"database/sql"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
	"gopkg.in/redis.v4"
)

const (
	dbDateFormat     = "2006-01-02"
	dbDatetimeFormat = "2006-01-02 15:04:05"
)

type contextKeyRedis struct{}

func OpenDB(dsn string, maxConnections int, logging bool) (*gorm.DB, error) {
	db, err := sql.Open(
		"mysql",
		dsn+"?charset=utf8mb4&parseTime=true&loc=UTC",
	)
	if err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to sql.Open()"),
		)
	}

	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxConnections)
	db.SetConnMaxLifetime(30 * time.Second)

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to gorm.Open()"),
		)
	}
	gormDB.LogMode(logging)

	return gormDB, nil
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

// TODO: Remove this function and use context_data
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
		return errors.NewInternalError(
			errors.WithError(tx.Error),
			errors.WithMessagef("Failed to begin transaction: name=%v", name),
		)
	}

	var err error
	success := false
	defer func() {
		if success {
			return
		}
		if err2 := tx.Rollback().Error; err2 != nil {
			err = errors.NewInternalError(
				errors.WithError(err2),
				errors.WithMessagef("Failed to rollback transaction: name=%v", name),
			)
		}
	}()

	if err2 := callback(tx); err2 != nil {
		return err2
	}
	if tx.Error != nil {
		return tx.Error
	}
	if err2 := tx.Commit().Error; err2 != nil {
		return errors.NewInternalError(
			errors.WithError(err2),
			errors.WithMessagef("Failed to commit transaction: name=%v", name),
		)
	}
	success = true
	return err
}

func LoadAllTables(db *gorm.DB, dbName string) ([]string, error) {
	type Table struct {
		Name string `gorm:"column:table_name"`
	}
	tables := []Table{}
	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = ?"
	if err := db.Raw(sql, dbName).Scan(&tables).Error; err != nil {
		return nil, err
	}

	tableNames := []string{}
	for _, t := range tables {
		if t.Name == "goose_db_version" {
			continue
		}
		tableNames = append(tableNames, t.Name)
	}
	return tableNames, nil
}

func ReplaceToTestDBURL(dbURL string) string {
	if strings.HasSuffix(dbURL, "/lekcije") {
		return strings.Replace(dbURL, "/lekcije", "/lekcije_test", 1)
	}
	return dbURL
}

func GetDBName(dbURL string) string {
	if index := strings.LastIndex(dbURL, "/"); index != -1 {
		return dbURL[index+1:]
	}
	return ""
}

func Placeholders(values []interface{}) string {
	s := strings.Repeat("?,", len(values))
	return strings.TrimRight(s, ",")
}

func wrapNotFound(result *gorm.DB, kind, key, value string) *errors.AnnotatedError {
	if result.RecordNotFound() {
		return errors.NewAnnotatedError(
			errors.CodeNotFound,
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(kind, key, value)),
		)
	} else {
		return nil
	}
}
