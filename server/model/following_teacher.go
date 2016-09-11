package model

import (
	"time"

	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type FollowingTeacher struct {
	UserId    uint32 `gorm:"primary_key"`
	TeacherId uint32 `gorm:"primary_key"`
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

func (s *FollowingTeacherService) FindTeachersByUserId(userId uint32) ([]*Teacher, error) {
	limit := 100
	values := make([]*Teacher, 0, limit)
	sql := fmt.Sprintf(`
	SELECT t.* FROM teacher AS t
	INNER JOIN following_teacher AS ft ON t.id = ft.teacher_id
	WHERE ft.user_id = ?
	ORDER BY ft.updated_at DESC, t.id ASC
	LIMIT %d
	`, limit) // TODO: OFFSET

	if result := s.db.Raw(strings.TrimSpace(sql), userId).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return values, nil
		}
		return values, errors.InternalWrapf(result.Error, "")
	}
	return values, nil
}

func (s *FollowingTeacherService) FindTeacherIdsByUserId(userId uint32) ([]uint32, error) {
	values := make([]*FollowingTeacher, 0, 1000)
	sql := `SELECT teacher_id FROM following_teacher WHERE user_id = ?`
	if result := s.db.Raw(sql, userId).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.InternalWrapf(result.Error, "")
	}
	ids := make([]uint32, len(values))
	for i, t := range values {
		ids[i] = t.TeacherId
	}
	return ids, nil
}

func (s *FollowingTeacherService) DeleteTeachersByUserIdAndTeacherIds(
	userId uint32,
	teacherIds []uint32,
) (int, error) {
	if len(teacherIds) == 0 {
		return 0, nil
	}

	placeholder := strings.TrimRight(strings.Repeat("?,", len(teacherIds)), ",")
	sql := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND teacher_id IN (%s)", s.TableName(), placeholder)
	args := []interface{}{userId}
	for _, teacherId := range teacherIds {
		args = append(args, teacherId)
	}
	if result := s.db.Exec(sql, args...); result.Error != nil {
		return 0, result.Error
	} else {
		return int(result.RowsAffected), nil
	}
}
