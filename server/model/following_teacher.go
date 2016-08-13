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

type FollowingTeacherRepository struct {
	db *gorm.DB
}

var FollowingTeacherRepo FollowingTeacherRepository

func (r *FollowingTeacherRepository) TableName() string {
	return (&FollowingTeacher{}).TableName()
}

func (r *FollowingTeacherRepository) FindTeachersByUserId(userId uint32) ([]*Teacher, error) {
	limit := 10
	values := make([]*Teacher, 0, limit)
	sql := fmt.Sprintf(`
	SELECT t.* FROM teacher AS t
	INNER JOIN following_teacher AS ft ON t.id = ft.teacher_id
	WHERE ft.user_id = ?
	ORDER BY ft.updated_at DESC
	LIMIT %d
	`, limit)

	if result := r.db.Raw(strings.TrimSpace(sql), userId).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return values, nil
		}
		return values, errors.InternalWrapf(result.Error, "")
	}
	return values, nil
}

func (r *FollowingTeacherRepository) DeleteTeachersByUserIdAndTeacherIds(
	userId uint32,
	teacherIds []uint32,
) (int, error) {
	if len(teacherIds) == 0 {
		return 0, nil
	}

	placeholder := strings.TrimRight(strings.Repeat("?,", len(teacherIds)), ",")
	sql := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND teacher_id IN (%s)", r.TableName(), placeholder)
	args := []interface{}{userId}
	for _, teacherId := range teacherIds {
		args = append(args, teacherId)
	}
	if result := r.db.Exec(sql, args...); result.Error != nil {
		return 0, result.Error
	} else {
		return int(result.RowsAffected), nil
	}
}
