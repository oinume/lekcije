package di

import (
	"database/sql"

	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"

	iga_measurement "github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	irollbar "github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/usecase"
)

func NewErrorRecorderUsecase(appLogger *zap.Logger, rollbarClient *rollbar.Client) *usecase.ErrorRecorder {
	return usecase.NewErrorRecorder(
		appLogger,
		irollbar.NewErrorRecorderRepository(rollbarClient),
	)
}

func NewFollowingTeacherUsecase(db *sql.DB) *usecase.FollowingTeacher {
	return usecase.NewFollowingTeacher(
		mysql.NewDBRepository(db),
		mysql.NewFollowingTeacherRepository(db),
	)
}

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
