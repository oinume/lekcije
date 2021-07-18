package mysqltest

import (
	"database/sql"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/repository"
)

type Repositories struct {
	db           repository.DB
	user         repository.User
	userAPIToken repository.UserAPIToken
	userGoogle   repository.UserGoogle
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		db:           mysql.NewDBRepository(db),
		user:         mysql.NewUserRepository(db),
		userAPIToken: mysql.NewUserAPITokenRepository(db),
		userGoogle:   mysql.NewUserGoogleRepository(db),
	}
}

func (r *Repositories) DB() repository.DB {
	return r.db
}

func (r *Repositories) User() repository.User {
	return r.user
}

func (r *Repositories) UserAPIToken() repository.UserAPIToken {
	return r.userAPIToken
}

func (r *Repositories) UserGoogle() repository.UserGoogle {
	return r.userGoogle
}
