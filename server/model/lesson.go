package model

import (
	"fmt"
	"time"

	"github.com/oinume/goenum"
)

type Lesson struct {
	TeacherId uint32
	Datetime  time.Time
	Status    string // TODO: enum
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
	Finished   int `goenum:"終了"`
	Reserved   int `goenum:"予約済"`
	Reservable int `goenum:"予約可"`
	Cancelled  int `goenum:"休講"`
}

var LessonStatuses = goenum.EnumerateStruct(&LessonStatus{
	Finished:   1,
	Reserved:   2,
	Reservable: 3,
	Cancelled:  4,
})
