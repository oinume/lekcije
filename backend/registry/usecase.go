package registry

import (
	"database/sql"

	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	iga_measurement "github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	irollbar "github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func NewErrorRecorderUsecase(appLogger *zap.Logger, rollbarClient *rollbar.Client) *usecase.ErrorRecorder {
	return usecase.NewErrorRecorder(
		appLogger,
		irollbar.NewErrorRecorderRepository(rollbarClient),
	)
}

func NewFollowingTeacherUsecase(appLogger *zap.Logger, db *sql.DB, mCountryList *model2.MCountryList) *usecase.FollowingTeacher {
	return usecase.NewFollowingTeacher(
		appLogger,
		mysql.NewDBRepository(db),
		mysql.NewFollowingTeacherRepository(db),
		mysql.NewUserRepository(db),
		mysql.NewTeacherRepository(db),
		dmm_eikaiwa.NewLessonFetcher(nil, 1, false, mCountryList, appLogger),
	)
}

func NewGAMeasurementUsecase(client iga_measurement.Client) *usecase.GAMeasurement {
	return usecase.NewGAMeasurement(iga_measurement.NewGAMeasurementRepository(client))
}

func NewLessonUsecase(db *sql.DB) *usecase.Lesson {
	return usecase.NewLesson(
		mysql.NewLessonRepository(db),
		mysql.NewLessonStatusLogRepository(db),
	)
}

func NewNotificationTimeSpanUsecase(db *sql.DB) *usecase.NotificationTimeSpan {
	return usecase.NewNotificationTimeSpan(mysql.NewNotificationTimeSpanRepository(db))
}

func NewStatNotifierUsecase(db *sql.DB) *usecase.StatNotifier {
	return usecase.NewStatNotifier(mysql.NewStatNotifierRepository(db))
}

func NewTeacherUsecase(db *sql.DB) *usecase.Teacher {
	return usecase.NewTeacher(mysql.NewTeacherRepository(db))
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
