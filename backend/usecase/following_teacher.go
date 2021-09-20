package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type FollowingTeacher struct {
	dbRepo repository.DB
}

func NewFollowingTeacher(dbRepo repository.DB) *FollowingTeacher {
	return &FollowingTeacher{
		dbRepo: dbRepo,
	}
}

func (u *FollowingTeacher) ListTeachersByUserID(ctx context.Context, userID uint32) ([]*model2.Teacher, error) {
	return nil, nil
}
