package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type Teacher struct {
	teacherRepo repository.Teacher
}

func NewTeacher(teacherRepo repository.Teacher) *Teacher {
	return &Teacher{
		teacherRepo: teacherRepo,
	}
}

func (u *Teacher) CreateOrUpdate(ctx context.Context, teacher *model2.Teacher) error {
	return u.teacherRepo.CreateOrUpdate(ctx, teacher)
}
