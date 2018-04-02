package model

import (
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
	user := helper.CreateRandomUser()
	count, err := followingTeacherService.CountFollowingTeachersByUserID(user.ID)
	a.Nil(err)
	a.Equal(0, count)
}

func TestFollowingTeacherService_FindTeacherIDsByUserID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser()
	teacher := helper.CreateRandomTeacher()
	now := time.Now()
	_, err := followingTeacherService.FollowTeacher(user.ID, teacher, now)
	r.NoError(err)

	teacherIDs, err := followingTeacherService.FindTeacherIDsByUserID(user.ID, 5)
	r.NoError(err)
	a.Equal(1, len(teacherIDs))

	err = teacherService.IncrementFetchErrorCount(teacher.ID, 6)
	r.NoError(err)
	teacherIDs, err = followingTeacherService.FindTeacherIDsByUserID(user.ID, 5)
	r.NoError(err)
	a.Equal(0, len(teacherIDs))
}
