package model

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestFollowingTeacherService_FollowTeacher(t *testing.T) {
	a := assert.New(t)

	user := createTestUser()
	teacher := &Teacher{
		ID:   1,
		Name: "Donald",
	}
	err := followingTeacherService.FollowTeacher(user.ID, teacher, time.Now().UTC())
	a.Nil(err)

	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
	a.Nil(err)
	a.Equal(1, len(teachers))
	a.Equal("Donald", teachers[0].Name)
}

func TestFollowingTeacherService_CountFollowingTeachersByUserID(t *testing.T) {
	a := assert.New(t)

	user := createTestUser()
	count, err := followingTeacherService.CountFollowingTeachersByUserID(user.ID)
	a.Nil(err)
	a.Equal(0, count)
}
