package model

import (
	"context"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFollowingTeacherService_FollowTeacher(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser()
	teacher := helper.CreateTeacher(1, "Donald")
	_, err := followingTeacherService.FollowTeacher(user.ID, teacher, time.Now().UTC())
	r.NoError(err)

	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
	r.NoError(err)
	a.Equal(1, len(teachers))
	a.Equal("Donald", teachers[0].Name)
}

func TestFollowingTeacherService_CountFollowingTeachersByUserID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser()
	count, err := followingTeacherService.CountFollowingTeachersByUserID(user.ID)
	r.Nil(err)
	a.Equal(0, count)
}

func TestFollowingTeacherService_FindTeacherIDsByUserID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	t.Run("Can find teacherIDs", func(t *testing.T) {
		user, teacher, err := followTeacher()
		r.NoError(err)

		teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(
			context.Background(),
			user.ID,
			5,
			defaultLastLessonAt,
		)
		r.NoError(err)
		a.Equal(1, len(teacherIDs))
		a.Equal(teacher.ID, teacherIDs[0])
	})

	t.Run("Cannot find teacherIDs because of fetchErrorCount", func(t *testing.T) {
		user, teacher, err := followTeacher()
		r.NoError(err)

		err = teacherService.IncrementFetchErrorCount(teacher.ID, 6)
		r.NoError(err)
		teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(
			context.Background(),
			user.ID,
			5,
			defaultLastLessonAt,
		)
		r.NoError(err)
		a.Equal(0, len(teacherIDs))
	})

	t.Run("Cannot find teacherIDs because of lastLessonAt", func(t *testing.T) {
		user, teacher, err := followTeacher()
		r.NoError(err)

		teacher.LastLessonAt = time.Now().Add(-1 * 3 * 24 * time.Hour)
		err = teacherService.CreateOrUpdate(teacher)
		r.NoError(err)
		teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(
			context.Background(),
			user.ID,
			5,
			time.Now(),
		)
		r.NoError(err)
		a.Equal(0, len(teacherIDs))
	})
}

func followTeacher() (*User, *Teacher, error) {
	user := helper.CreateRandomUser()
	teacher := helper.CreateRandomTeacher()
	now := time.Now()
	_, err := followingTeacherService.FollowTeacher(user.ID, teacher, now)
	if err != nil {
		return nil, nil, err
	}
	return user, teacher, nil
}
