package repository

import (
	"context"
	"time"

	"github.com/oinume/lekcije/backend/model2"
)

type FollowingTeacher interface {
	CountFollowingTeachersByUserID(ctx context.Context, userID uint) (int, error)
	Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error
	DeleteByUserIDAndTeacherIDs(ctx context.Context, userID uint, teacherIDs []uint) error
	FindTeacherIDsByUserID(ctx context.Context, userID uint, fetchErrorCount int, lastLessonAt time.Time) ([]uint, error)
	FindTeachersByUserID(ctx context.Context, userID uint) ([]*model2.Teacher, error)
	FindByUserID(ctx context.Context, userID uint) ([]*model2.FollowingTeacher, error)
	FindByUserIDAndTeacherID(ctx context.Context, userID uint, teacherID uint) (*model2.FollowingTeacher, error)
}
