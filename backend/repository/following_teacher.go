package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type FollowingTeacher interface {
	CountFollowingTeachersByUserID(ctx context.Context, userID uint) (int, error)
	Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error
	FindTeachersByUserID(ctx context.Context, userID uint) ([]*model2.Teacher, error)
}
