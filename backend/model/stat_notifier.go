package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
)

type StatNotifier struct {
	Datetime             time.Time
	Interval             uint8
	Elapsed              uint32
	UserCount            uint32
	FollowedTeacherCount uint32
}

func (*StatNotifier) TableName() string {
	return "stat_notifier"
}

type StatNotifierService struct {
	db *gorm.DB
}

func NewStatNotifierService(db *gorm.DB) *StatNotifierService {
	return &StatNotifierService{db}
}

func (s *StatNotifierService) CreateOrUpdate(v *StatNotifier) error {
	datetime := v.Datetime.Format(dbDatetimeFormat)
	sql := fmt.Sprintf(`INSERT INTO %s VALUES (?, ?, ?, ?, ?)`, v.TableName())
	sql += " ON DUPLICATE KEY UPDATE"
	sql += " elapsed=?, user_count=?, followed_teacher_count=?"
	values := []interface{}{
		datetime, v.Interval, v.Elapsed, v.UserCount, v.FollowedTeacherCount,
		v.Elapsed, v.UserCount, v.FollowedTeacherCount,
	}
	if err := s.db.Exec(sql, values...).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource(v.TableName(), "datetime", datetime)),
		)
	}
	return nil
}
