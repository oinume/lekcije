package model2

import (
	"github.com/volatiletech/null/v8"

	"github.com/oinume/lekcije/backend/model"
)

func NewUserFromModel(u *model.User) *User {
	var emailVerified uint8
	if u.EmailVerified {
		emailVerified = 1
	}
	return &User{
		ID:                 uint(u.ID),
		Name:               u.Name,
		Email:              u.Email,
		EmailVerified:      emailVerified,
		PlanID:             u.PlanID,
		FollowedTeacherAt:  null.NewTime(u.FollowedTeacherAt.Time, u.FollowedTeacherAt.Valid),
		OpenNotificationAt: null.NewTime(u.OpenNotificationAt.Time, u.OpenNotificationAt.Valid),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}
