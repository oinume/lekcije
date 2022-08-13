package usecase

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/model2"
)

const MaxFollowTeacherCount = 20

type FollowingTeacher struct {
	appLogger            *zap.Logger
	dbRepo               repository.DB
	followingTeacherRepo repository.FollowingTeacher
	mCountryRepo         repository.MCountry
	userRepo             repository.User
	teacherRepo          repository.Teacher
}

func NewFollowingTeacher(
	appLogger *zap.Logger,
	dbRepo repository.DB,
	followingTeacherRepo repository.FollowingTeacher,
	mCountryRepo repository.MCountry,
	userRepo repository.User,
	teacherRepo repository.Teacher,
) *FollowingTeacher {
	return &FollowingTeacher{
		appLogger:            appLogger,
		dbRepo:               dbRepo,
		followingTeacherRepo: followingTeacherRepo,
		mCountryRepo:         mCountryRepo,
		userRepo:             userRepo,
		teacherRepo:          teacherRepo,
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

func (u *FollowingTeacher) FollowTeacher(ctx context.Context, user *model2.User, teacher *model2.Teacher) (bool, error) {
	reachesLimit, err := u.ReachesFollowingTeacherLimit(ctx, user.ID, 1)
	if err != nil {
		return false, err
	}
	if reachesLimit {
		return false, errors.NewFailedPreconditionError(errors.WithMessagef("フォロー可能な上限数(%d)を超えました", MaxFollowTeacherCount))
	}

	// Update user.followed_teacher_at when first time to follow teachers.
	// the column is used for showing tutorial or not.
	updateFollowedTeacherAt := false
	if !user.FollowedTeacherAt.Valid {
		now := time.Now().UTC()
		if err := u.userRepo.UpdateFollowedTeacherAt(ctx, user.ID, now); err != nil {
			return false, err
		}
		if err := u.userRepo.UpdateOpenNotificationAt(ctx, user.ID, now); err != nil {
			return false, err
		}
		updateFollowedTeacherAt = true
	}

	mCountries, err := u.mCountryRepo.FindAll(ctx)
	if err != nil {
		return false, err
	}
	// TODO: DI
	fetcher := dmm_eikaiwa.NewLessonFetcher(nil, 1, false, model2.NewMCountryList(mCountries), u.appLogger)
	defer fetcher.Close()
	fetchedTeacher, _, err := fetcher.Fetch(ctx, teacher.ID)
	if err != nil {
		return false, err
	}

	if err := u.teacherRepo.CreateOrUpdate(ctx, fetchedTeacher); err != nil {
		return false, err
	}
	if err := u.followingTeacherRepo.Create(ctx, &model2.FollowingTeacher{
		UserID:    user.ID,
		TeacherID: teacher.ID,
	}); err != nil {
		return false, err
	}
	return updateFollowedTeacherAt, nil
}

func (u *FollowingTeacher) ReachesFollowingTeacherLimit(ctx context.Context, userID uint, additionalTeachers int) (bool, error) {
	count, err := u.followingTeacherRepo.CountFollowingTeachersByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return count+additionalTeachers > MaxFollowTeacherCount, nil // TODO: Refer plan's limit
}
