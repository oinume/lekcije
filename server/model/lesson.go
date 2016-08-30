package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/goenum"
	"github.com/oinume/lekcije/server/errors"
)

type Lesson struct {
	TeacherId uint32    `gorm:"primary_key"`
	Datetime  time.Time `gorm:"primary_key"`
	Status    string    // TODO: enum
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherId: %v, Datetime: %v, Status: %v",
		l.TeacherId, l.Datetime, l.Status,
	)
}

type LessonStatus struct {
	Finished  int `goenum:"終了"`
	Reserved  int `goenum:"予約済"`
	Available int `goenum:"予約可"`
	Cancelled int `goenum:"休講"`
}

var LessonStatuses = goenum.EnumerateStruct(&LessonStatus{
	Finished:  1,
	Reserved:  2,
	Available: 3,
	Cancelled: 4,
})

type LessonServiceType struct {
	db *gorm.DB
}

var LessonService LessonServiceType

func (s *LessonServiceType) TableName() string {
	return (&Lesson{}).TableName()
}

func (s *LessonServiceType) UpdateLessons(lessons []*Lesson) (int64, error) {
	if len(lessons) == 0 {
		return 0, nil
	}

	now := time.Now()
	sql := fmt.Sprintf("INSERT INTO %s VALUES", s.TableName())
	values := []interface{}{}
	for _, l := range lessons {
		sql += " (?, ?, ?, ?, ?),"
		values = append(values, l.TeacherId, l.Datetime, l.Status, now.Format(dbDatetimeFormat)) // TODO: enum?
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += " ON DUPLICATE KEY UPDATE status=VALUES(status), updated_at=VALUES(updated_at)"

	result := s.db.Exec(sql, values)
	if err := result.Error; err != nil {
		return 0, errors.InternalWrapf(err, "")
	}

	return result.RowsAffected, nil
}
