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
	teacherID := uint32(1)
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(teacherID, datetime, "Reserved", 5)

	affected, err := lessonService.UpdateLessons(lessons)
	a.Nil(err)
	a.Equal(int64(5), affected)

	foundLessons, err := lessonService.FindLessons(teacherID, datetime, datetime)
	a.Nil(err)
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
	affected, err := lessonService.UpdateLessons(lessons)
	a.Nil(err)
	a.Equal(int64(2), affected) // Why 2????
	foundLessons, err := lessonService.FindLessons(1, datetime, datetime)
	a.Nil(err)
	a.Equal(strings.ToLower(foundLessons[0].Status), "available")
}

func TestGetNewAvailableLessons1(t *testing.T) {
	a := assert.New(t)

	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons1 := createLessons(1, datetime, "Reserved", 3)
	lessons2 := createLessons(1, datetime, "Reserved", 3)
	lessons2[1].Status = "Available"
	// Test GetNewAvailableLessons returns a lesson when new lesson is "Available"
	availableLessons := lessonService.GetNewAvailableLessons(lessons1, lessons2)
	a.Equal(1, len(availableLessons))
	a.Equal(datetime.Add(1*time.Hour), availableLessons[0].Datetime)
}

func TestGetNewAvailableLessons2(t *testing.T) {
	a := assert.New(t)

	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons1 := createLessons(1, datetime, "Reserved", 3)
	lessons2 := createLessons(1, datetime, "Reserved", 3)
	lessons1[0].Status = "Available"
	lessons2[0].Status = "Available"
	// Test GetNewAvailableLessons returns nothing when both lessons are "Available"
	availableLessons := lessonService.GetNewAvailableLessons(lessons1, lessons2)
	a.Equal(0, len(availableLessons))
}

func createLessons(teacherID uint32, baseDatetime time.Time, status string, length int) []*Lesson {
	lessons := make([]*Lesson, length)
	now := time.Now()
	for i := range lessons {
		lessons[i] = &Lesson{
			TeacherID: teacherID,
			Datetime:  baseDatetime.Add(time.Duration(i) * time.Hour),
			Status:    status,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return lessons
}
