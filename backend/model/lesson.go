package model

import (
	"fmt"
	"time"

	"github.com/oinume/goenum"
)

const (
	lessonTimeFormat = "2006-01-02 15:04"
)

type Lesson struct {
	ID        uint64 `gorm:"primary_key"`
	TeacherID uint32
	Datetime  time.Time
	Status    string // TODO: enum
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherID=%v, Datetime=%v, Status=%v",
		l.TeacherID, l.Datetime.Format(lessonTimeFormat), l.Status,
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
