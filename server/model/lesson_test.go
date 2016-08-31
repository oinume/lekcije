package model

import (
	"strings"
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
			Status:    "Reserved",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	affected, err := LessonService.UpdateLessons(lessons)
	a.NoError(err)
	a.Equal(int64(5), affected)

	foundLessons, err := LessonService.FindLessons(1, datetime, datetime)
	a.NoError(err)
	a.Equal(len(lessons), len(foundLessons))
	for i := range lessons {
		// TODO: custom enum type
		a.Equal(strings.ToLower(lessons[i].Status), strings.ToLower(foundLessons[i].Status))
	}
}
