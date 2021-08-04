package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/backend/randoms"
)

func TestNewTeachersFromIDOrURL(t *testing.T) {
	tests := map[string]struct {
		ids       string
		err       error
		nTeachers int
	}{
		"normal": {
			ids:       "1,2",
			err:       nil,
			nTeachers: 2,
		},
		"trailing_comma": {
			ids:       "1,2,3,",
			err:       nil,
			nTeachers: 3,
		},
		"empty": {
			ids:       "",
			err:       fmt.Errorf("code.InvalidArgument: Failed to parse idsOrURL: "),
			nTeachers: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			teachers, err := NewTeachersFromIDsOrURL(test.ids)
			if test.err == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.err != nil {
				if got, want := err.Error(), test.err.Error(); got != want {
					t.Fatalf("unexpected error: got=%q, want=%q", got, want)
				}
			}
			if got, want := len(teachers), test.nTeachers; got != want {
				t.Errorf("unexpected teachers length: got=%v, want=%v", got, want)
			}
		})
	}
}

func TestTeacherService_CreateOrUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	teacher := &Teacher{
		ID:                uint32(randoms.MustNewInt64(9999999)),
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
		ID:                uint32(randoms.MustNewInt64(9999999)),
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
