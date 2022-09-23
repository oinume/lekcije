package usecase

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
)

const MaxFollowTeacherCount = 20

type FollowingTeacher struct {
	appLogger            *zap.Logger
	dbRepo               repository.DB
	followingTeacherRepo repository.FollowingTeacher
	userRepo             repository.User
	teacherRepo          repository.Teacher
	lessonFetcherRepo    repository.LessonFetcher
}

func NewFollowingTeacher(
	appLogger *zap.Logger,
	dbRepo repository.DB,
	followingTeacherRepo repository.FollowingTeacher,
	userRepo repository.User,
	teacherRepo repository.Teacher,
	lessonFetcherRepo repository.LessonFetcher,
) *FollowingTeacher {
	return &FollowingTeacher{
		appLogger:            appLogger,
		dbRepo:               dbRepo,
		followingTeacherRepo: followingTeacherRepo,
		userRepo:             userRepo,
		teacherRepo:          teacherRepo,
		lessonFetcherRepo:    lessonFetcherRepo,
	}
}

func (u *FollowingTeacher) Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error {
	return u.followingTeacherRepo.Create(ctx, followingTeacher)
}

func (u *FollowingTeacher) DeleteFollowingTeachers(ctx context.Context, userID uint, teacherIDs []uint) error {
	return u.followingTeacherRepo.DeleteByUserIDAndTeacherIDs(ctx, userID, teacherIDs)
}

func (u *FollowingTeacher) FindTeachersByUserID(ctx context.Context, userID uint) ([]*model2.Teacher, error) {
	return u.followingTeacherRepo.FindTeachersByUserID(ctx, userID)
}

func (u *FollowingTeacher) FollowTeacher(ctx context.Context, user *model2.User, teacher *model2.Teacher) (*model2.FollowingTeacher, bool, error) {
	reachesLimit, err := u.ReachesFollowingTeacherLimit(ctx, user.ID, 1)
	if err != nil {
		return nil, false, err
	}
	if reachesLimit {
		return nil, false, errors.NewFailedPreconditionError(errors.WithMessagef("フォロー可能な上限数(%d)を超えました", MaxFollowTeacherCount))
	}

	// Update user.followed_teacher_at when first time to follow teachers.
	// the column is used for showing tutorial or not.
	updateFollowedTeacherAt := false
	if !user.FollowedTeacherAt.Valid {
		now := time.Now().UTC()
		if err := u.userRepo.UpdateFollowedTeacherAt(ctx, user.ID, now); err != nil {
			return nil, false, err
		}
		if err := u.userRepo.UpdateOpenNotificationAt(ctx, user.ID, now); err != nil {
			return nil, false, err
		}
		updateFollowedTeacherAt = true
	}

	// TODO: Close
	fetchedTeacher, _, err := u.lessonFetcherRepo.Fetch(ctx, teacher.ID)
	if err != nil {
		return nil, false, err
	}

	if err := u.teacherRepo.CreateOrUpdate(ctx, fetchedTeacher); err != nil {
		return nil, false, err
	}
	ft := &model2.FollowingTeacher{
		UserID:    user.ID,
		TeacherID: teacher.ID,
	}
	if err := u.followingTeacherRepo.Create(ctx, ft); err != nil {
		return nil, false, err
	}
	return ft, updateFollowedTeacherAt, nil
}

func (u *FollowingTeacher) ReachesFollowingTeacherLimit(ctx context.Context, userID uint, additionalTeachers int) (bool, error) {
	count, err := u.followingTeacherRepo.CountFollowingTeachersByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return count+additionalTeachers > MaxFollowTeacherCount, nil // TODO: Refer plan's limit
}
