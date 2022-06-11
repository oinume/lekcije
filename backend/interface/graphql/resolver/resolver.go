package resolver

import "github.com/oinume/lekcije/backend/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	followingTeacherRepo     repository.FollowingTeacher
	notificationTimeSpanRepo repository.NotificationTimeSpan
	userRepo                 repository.User
	teacherRepo              repository.Teacher
}

func NewResolver(
	followingTeacherRepo repository.FollowingTeacher,
	notificationTimeSpanRepo repository.NotificationTimeSpan,
	teacherRepo repository.Teacher,
	userRepo repository.User,
) *Resolver {
	return &Resolver{
		followingTeacherRepo:     followingTeacherRepo,
		notificationTimeSpanRepo: notificationTimeSpanRepo,
		teacherRepo:              teacherRepo,
		userRepo:                 userRepo,
	}
}
