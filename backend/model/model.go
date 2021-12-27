package model

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/oinume/lekcije/backend/errors"
)

const (
	dbDateFormat     = "2006-01-02"
	dbDatetimeFormat = "2006-01-02 15:04:05"
)

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

func ReplaceToTestDBURL(t *testing.T, dbURL string) string {
	if strings.HasSuffix(dbURL, "/lekcije") {
		return strings.Replace(dbURL, "/lekcije", "/lekcije_test", 1)
	}
	return dbURL
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
