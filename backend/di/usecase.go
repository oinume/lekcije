package di

import (
	"database/sql"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/usecase"
)

func NewUserUsecase(db *sql.DB) *usecase.User {
	return usecase.NewUser(
		mysql.NewDBRepository(db),
		mysql.NewUserRepository(db),
		mysql.NewUserGoogleRepository(db),
	)
}

func NewUserAPITokenUsecase(db *sql.DB) *usecase.UserAPIToken {
	return usecase.NewUserAPIToken(
		mysql.NewUserAPITokenRepository(db),
	)
}
