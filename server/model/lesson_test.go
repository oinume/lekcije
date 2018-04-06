package model

import (
	"strings"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLessonService_UpdateLessons(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacherID := uint32(util.RandomInt(999999))
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(teacherID, datetime, "Reserved", 5)

	affected, err := lessonService.UpdateLessons(lessons)
	r.NoError(err)
	a.Equal(int64(5), affected)
	for _, l := range lessons {
		a.NotEqual(uint64(0), l.ID)
		logs, err := lessonStatusLogService.FindAllByLessonID(l.ID)
		r.NoError(err)
		a.Equal(1, len(logs))
	}

	foundLessons, err := lessonService.FindLessons(teacherID, datetime, datetime)
	r.NoError(err)
	a.Equal(len(lessons), len(foundLessons))
	for i := range lessons {
		// TODO: custom enum type
		a.Equal(strings.ToLower(lessons[i].Status), strings.ToLower(foundLessons[i].Status))
	}
}

func TestLessonService_UpdateLessonsOverwrite(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacherID := uint32(util.RandomInt(999999))
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(teacherID, datetime, "Available", 5)
	affected, err := lessonService.UpdateLessons(lessons)
	r.NoError(err)
	a.EqualValues(len(lessons), affected)

	time.Sleep(1 * time.Second)
	lessons[0].Status = "Reserved"
	affected, err = lessonService.UpdateLessons(lessons)
	r.NoError(err)
	a.EqualValues(1, affected)

	foundLessons, err := lessonService.FindLessons(teacherID, datetime, datetime)
	r.NoError(err)
	a.Equal(strings.ToLower(foundLessons[0].Status), "reserved")

	logs, err := lessonStatusLogService.FindAllByLessonID(foundLessons[0].ID)
	r.NoError(err)
	a.Equal(2, len(logs))
}

func TestLessonService_UpdateLessonsNoChange(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacherID := uint32(util.RandomInt(999999))
	datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalTimezone())
	lessons := createLessons(teacherID, datetime, "Available", 5)
	affected, err := lessonService.UpdateLessons(lessons)
	r.NoError(err)
	a.EqualValues(len(lessons), affected)

	affected, err = lessonService.UpdateLessons(lessons)
	r.NoError(err)
	a.EqualValues(0, affected)

	foundLessons, err := lessonService.FindLessons(teacherID, datetime, datetime)
	r.NoError(err)
	a.Equal(strings.ToLower(foundLessons[0].Status), "available")

	logs, err := lessonStatusLogService.FindAllByLessonID(foundLessons[0].ID)
	r.NoError(err)
	a.Equal(1, len(logs))
}

func TestLessonService_GetNewAvailableLessons1(t *testing.T) {
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

func TestLessonService_GetNewAvailableLessons2(t *testing.T) {
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
	now := time.Now().UTC()
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
