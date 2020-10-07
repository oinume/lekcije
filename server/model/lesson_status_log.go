package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/server/errors"
)

type LessonStatusLog struct {
	LessonID  uint64
	Status    string
	CreatedAt time.Time
}

func (*LessonStatusLog) TableName() string {
	return "lesson_status_log"
}

type LessonStatusLogService struct {
	db *gorm.DB
}

func NewLessonStatusLogService(db *gorm.DB) *LessonStatusLogService {
	return &LessonStatusLogService{db: db}
}

func (s *LessonStatusLogService) TableName() string {
	return (&LessonStatusLog{}).TableName()
}

func (s *LessonStatusLogService) Create(log *LessonStatusLog) error {
	return s.db.Create(log).Error
}

func (s *LessonStatusLogService) FindAllByLessonID(lessonID uint64) ([]*LessonStatusLog, error) {
	logs := make([]*LessonStatusLog, 0, 100)
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE lesson_id = ?`, s.TableName())
	result := s.db.Raw(sql, lessonID).Scan(&logs)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("FindAllByLessonID failed"),
			errors.WithResource(errors.NewResource(s.TableName(), "lessonID", lessonID)),
		)
	}
	return logs, nil
}
