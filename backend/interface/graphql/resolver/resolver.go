package resolver

import (
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	followingTeacherRepo        repository.FollowingTeacher
	followingTeacherUsecase     *usecase.FollowingTeacher
	notificationTimeSpanRepo    repository.NotificationTimeSpan
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan
	teacherRepo                 repository.Teacher
	userRepo                    repository.User
	userUsecase                 *usecase.User
}

func NewResolver(
	followingTeacherRepo repository.FollowingTeacher,
	followingTeacherUsecase *usecase.FollowingTeacher,
	notificationTimeSpanRepo repository.NotificationTimeSpan,
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan,
	teacherRepo repository.Teacher,
	userRepo repository.User,
	userUsecase *usecase.User,
) *Resolver {
	return &Resolver{
		followingTeacherRepo:        followingTeacherRepo,
		followingTeacherUsecase:     followingTeacherUsecase,
		notificationTimeSpanRepo:    notificationTimeSpanRepo,
		notificationTimeSpanUsecase: notificationTimeSpanUsecase,
		teacherRepo:                 teacherRepo,
		userRepo:                    userRepo,
		userUsecase:                 userUsecase,
	}
}
