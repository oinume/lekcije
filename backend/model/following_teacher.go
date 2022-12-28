package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
)

const MaxFollowTeacherCount = 20

type FollowingTeacher struct {
	UserID    uint32 `gorm:"primary_key"`
	TeacherID uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*FollowingTeacher) TableName() string {
	return "following_teacher"
}

type FollowingTeacherService struct {
	db *gorm.DB
}

func NewFollowingTeacherService(db *gorm.DB) *FollowingTeacherService {
	return &FollowingTeacherService{db: db}
}

func (s *FollowingTeacherService) TableName() string {
	return (&FollowingTeacher{}).TableName()
}

func (s *FollowingTeacherService) FindTeacherIDs() ([]uint32, error) {
	values := make([]*FollowingTeacher, 0, 5000)
	sql := `SELECT teacher_id FROM following_teacher LIMIT 5000`
	if result := s.db.Raw(sql).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessagef("failed to select teacher ids"),
		)
	}
	ids := make([]uint32, len(values))
	for i, t := range values {
		ids[i] = t.TeacherID
	}
	return ids, nil
}
