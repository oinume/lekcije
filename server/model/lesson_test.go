package model

import (
	"testing"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLessons(t *testing.T) {
	a := assert.New(t)

	lessons := make([]*Lesson, 5)
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	now := time.Now()
	for i := range lessons {
		lessons[i] = &Lesson{
			TeacherId: 1,
			Datetime:  datetime.Add(time.Duration(i) * time.Hour),
			Status:    "Available",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	affected, err := LessonService.UpdateLessons(lessons)
	a.NoError(err)
	a.Equal(int64(5), affected)
}
