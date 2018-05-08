package model

import (
	"testing"
	"time"

	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTeachersFromIDOrURL(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	teachers, err := NewTeachersFromIDsOrURL("1,2")
	r.Nil(err)
	a.Equal(2, len(teachers))

	teachers2, err := NewTeachersFromIDsOrURL("1,2,3,")
	r.Nil(err)
	a.Equal(3, len(teachers2))

	teachers3, err := NewTeachersFromIDsOrURL("")
	r.Error(err)
	a.Equal(0, len(teachers3))
}

func TestTeacherService_CreateOrUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacher := &Teacher{
		ID:                uint32(util.RandomInt(9999999)),
		Name:              "Donald",
		CountryID:         688, // Serbia
		Gender:            "male",
		Birthday:          time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		YearsOfExperience: 2,
		FavoriteCount:     100,
		ReviewCount:       50,
		Rating:            4.75,
		LastLessonAt:      time.Date(2018, 3, 1, 11, 10, 0, 0, time.UTC),
	}
	err := teacherService.CreateOrUpdate(teacher)
	r.NoError(err)

	actual, err := teacherService.FindByPK(teacher.ID)
	r.NoError(err)
	a.Equal(teacher.Name, actual.Name)
	a.Equal(teacher.Rating, actual.Rating)
	a.Equal(teacher.LastLessonAt, actual.LastLessonAt)

	newLastLessonAt := time.Date(2018, 4, 1, 11, 10, 0, 0, time.UTC)
	teacher.LastLessonAt = newLastLessonAt
	err = teacherService.CreateOrUpdate(teacher)
	r.NoError(err)
	actual, err = teacherService.FindByPK(teacher.ID)
	r.NoError(err)
	a.Equal(newLastLessonAt, actual.LastLessonAt)
}

func TestTeacherService_CreateOrUpdate2(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacher := &Teacher{
		ID:                uint32(util.RandomInt(9999999)),
		Name:              "Donald",
		CountryID:         688, // Serbia
		Gender:            "male",
		Birthday:          time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		YearsOfExperience: 2,
		FavoriteCount:     100,
		ReviewCount:       50,
		Rating:            5.0,
	}
	err := teacherService.CreateOrUpdate(teacher)
	r.NoError(err)

	actual, err := teacherService.FindByPK(teacher.ID)
	r.NoError(err)
	a.Equal(teacher.Name, actual.Name)
	a.Equal(defaultLastLessonAt, actual.LastLessonAt)
}

func TestTeacherService_IncrementFetchErrorCount(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacher := &Teacher{
		ID:     1,
		Name:   "test",
		Gender: "male",
	}
	err := teacherService.CreateOrUpdate(teacher)
	r.Nil(err)

	err = teacherService.IncrementFetchErrorCount(teacher.ID, 1)
	r.Nil(err)
	teacher2, err := teacherService.FindByPK(teacher.ID)
	r.Nil(err)
	a.Equal(uint8(1), teacher2.FetchErrorCount)
}
