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
	teacherId := uint32(1)
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(teacherId, datetime, "Reserved", 5)

	affected, err := LessonService.UpdateLessons(lessons)
	a.NoError(err)
	a.Equal(int64(5), affected)

	foundLessons, err := LessonService.FindLessons(teacherId, datetime, datetime)
	a.NoError(err)
	a.Equal(len(lessons), len(foundLessons))
	for i := range lessons {
		// TODO: custom enum type
		a.Equal(strings.ToLower(lessons[i].Status), strings.ToLower(foundLessons[i].Status))
	}
}

func TestUpdateLessonsOverwrite(t *testing.T) {
	a := assert.New(t)
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(1, datetime, "Reserved", 5)

	lessons[0].Status = "Available"
	affected, err := LessonService.UpdateLessons(lessons)
	a.NoError(err)
	a.Equal(int64(2), affected) // Why 2????
	foundLessons, err := LessonService.FindLessons(1, datetime, datetime)
	a.Equal(strings.ToLower(foundLessons[0].Status), "available")
}

func TestGetNewAvailableLessons(t *testing.T) {
	//a := assert.New(t)

	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons1 := createLessons(1, datetime, "Reserved", 3)
	lessons2 := createLessons(1, datetime, "Reserved", 3)
	lessons2[0].Status = "Available"

	LessonService.GetNewAvailableLessons(lessons1, lessons2)
}

func createLessons(teacherId uint32, baseDatetime time.Time, status string, length int) []*Lesson {
	lessons := make([]*Lesson, length)
	now := time.Now()
	for i := range lessons {
		lessons[i] = &Lesson{
			TeacherId: teacherId,
			Datetime:  baseDatetime.Add(time.Duration(i) * time.Hour),
			Status:    status,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return lessons
}
