package di

import (
	"database/sql"

	"github.com/oinume/lekcije/backend/ga_measurement"
	iga_measurement "github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/usecase"
)

func NewGAMeasurementUsecase(client ga_measurement.Client) *usecase.GAMeasurement {
	return usecase.NewGAMeasurement(iga_measurement.NewGAMeasurementRepository(client))
}

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
