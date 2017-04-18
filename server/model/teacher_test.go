package model

import (
	"testing"

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

func createTestTeacher(id uint32, name string) *Teacher {
	teacher := &Teacher{
		ID:   1,
		Name: "test",
	}
	if id != 0 {
		teacher.ID = id
	}
	if name != "" {
		teacher.Name = name
	}
	teacher.Gender = "female"
	if err := teacherService.CreateOrUpdate(teacher); err != nil {
		panic(err)
	}
	return teacher
}
