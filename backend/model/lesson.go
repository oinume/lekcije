package model

import (
	"github.com/oinume/goenum"
)

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
