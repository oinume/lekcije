package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type FollowingTeacher struct {
	dbRepo               repository.DB
	followingTeacherRepo repository.FollowingTeacher
}

func NewFollowingTeacher(dbRepo repository.DB, followingTeacherRepo repository.FollowingTeacher) *FollowingTeacher {
	return &FollowingTeacher{
		dbRepo:               dbRepo,
		followingTeacherRepo: followingTeacherRepo,
	}
}

func (u *FollowingTeacher) Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error {
	return u.followingTeacherRepo.Create(ctx, followingTeacher)
}

func (u *FollowingTeacher) FindTeachersByUserID(ctx context.Context, userID uint) ([]*model2.Teacher, error) {
	return u.followingTeacherRepo.FindTeachersByUserID(ctx, userID)
}
