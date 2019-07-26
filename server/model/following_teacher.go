package model

import (
	"context"
	"time"

	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
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

func (s *FollowingTeacherService) FindTeachersByUserID(userID uint32) ([]*Teacher, error) {
	limit := 100
	values := make([]*Teacher, 0, limit)
	sql := fmt.Sprintf(`
	SELECT t.* FROM teacher AS t
	INNER JOIN following_teacher AS ft ON t.id = ft.teacher_id
	WHERE ft.user_id = ?
	ORDER BY ft.updated_at DESC, t.id ASC
	LIMIT %d
	`, limit) // TODO: OFFSET

	if result := s.db.Raw(strings.TrimSpace(sql), userID).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return values, nil
		}
		return values, errors.NewInternalError(
			errors.WithError(result.Error),
		)
	}
	return values, nil
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

func (s *FollowingTeacherService) FindTeacherIDsByUserID(
	ctx context.Context, userID uint32, fetchErrorCount int, lastLessonAt time.Time,
) ([]uint32, error) {
	values := make([]*FollowingTeacher, 0, 1000)
	sql := `
	SELECT ft.teacher_id FROM following_teacher AS ft
	INNER JOIN teacher AS t ON ft.teacher_id = t.id
	WHERE
      user_id = ?
      AND t.fetch_error_count <= ?
      AND t.last_lesson_at >= ?
	`
	sql = strings.TrimSpace(sql)
	if result := s.db.Raw(sql, userID, fetchErrorCount, lastLessonAt.Format(dbDatetimeFormat)).Scan(&values); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
		)
	}
	ids := make([]uint32, len(values))
	for i, t := range values {
		ids[i] = t.TeacherID
	}
	return ids, nil
}

func (s *FollowingTeacherService) CountFollowingTeachersByUserID(userID uint32) (int, error) {
	count := struct {
		Count int
	}{}
	sql := `SELECT COUNT(*) AS count FROM following_teacher WHERE user_id = ?`
	if err := s.db.Raw(sql, userID).Scan(&count).Error; err != nil {
		return 0, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("count failed"),
			errors.WithResource(errors.NewResource("following_teacher", "userID", userID)),
		)
	}
	return count.Count, nil
}

func (s FollowingTeacherService) ReachesFollowingTeacherLimit(userID uint32, additionalTeachers int) (bool, error) {
	count, err := s.CountFollowingTeachersByUserID(userID)
	if err != nil {
		return false, err
	}
	return count+additionalTeachers > MaxFollowTeacherCount, nil // TODO: Refer plan's limit
}

func (s *FollowingTeacherService) FollowTeacher(
	userID uint32, teacher *Teacher, timestamp time.Time,
) (*FollowingTeacher, error) {
	// Create teacher at first
	teacher.CreatedAt = timestamp
	teacher.UpdatedAt = timestamp
	teacherService := NewTeacherService(s.db)
	if err := teacherService.CreateOrUpdate(teacher); err != nil {
		return nil, err
	}

	ft := &FollowingTeacher{
		UserID:    userID,
		TeacherID: teacher.ID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
	if err := s.db.FirstOrCreate(ft).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(
				errors.NewResourceWithEntries(ft.TableName(), []errors.ResourceEntry{
					{"userID", userID},
					{"teacherID", teacher.ID},
				}),
			),
		)
	}
	return ft, nil
}

func (s *FollowingTeacherService) DeleteTeachersByUserIDAndTeacherIDs(
	userID uint32,
	teacherIDs []uint32,
) (int, error) {
	if len(teacherIDs) == 0 {
		return 0, nil
	}

	placeholder := strings.TrimRight(strings.Repeat("?,", len(teacherIDs)), ",")
	sql := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND teacher_id IN (%s)", s.TableName(), placeholder)
	args := []interface{}{userID}
	for _, teacherID := range teacherIDs {
		args = append(args, teacherID)
	}
	if result := s.db.Exec(sql, args...); result.Error != nil {
		return 0, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to delete following_teacher"),
			errors.WithResource(errors.NewResourceWithEntries(s.TableName(), []errors.ResourceEntry{
				{"userID", userID},
				{"teacherIDs", teacherIDs},
			})),
		)
	} else {
		return int(result.RowsAffected), nil
	}
}
