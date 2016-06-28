package model

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

func Open() (*gorm.DB, error) {
	dbDsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%v?charset=utf8mb4&parseTime=true&loc=UTC", dbDsn),
	)
	return db, err
}
