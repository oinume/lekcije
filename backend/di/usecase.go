package di

import (
	"database/sql"

	iga_measurement "github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/usecase"
)

func NewGAMeasurementUsecase(client iga_measurement.Client) *usecase.GAMeasurement {
	return usecase.NewGAMeasurement(iga_measurement.NewGAMeasurementRepository(client))
}

func NewNotificationTimeSpanUsecase(db *sql.DB) *usecase.NotificationTimeSpan {
	return usecase.NewNotificationTimeSpan(mysql.NewNotificationTimeSpanRepository(db))
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
