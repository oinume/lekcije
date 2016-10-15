package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTeachersFromIDOrURL(t *testing.T) {
	a := assert.New(t)
	teachers, err := NewTeachersFromIDsOrURL("1,2")
	a.Nil(err)
	a.Equal(2, len(teachers))

	teachers2, err := NewTeachersFromIDsOrURL("1,2,3,")
	a.Nil(err)
	a.Equal(3, len(teachers2))

	teachers3, err := NewTeachersFromIDsOrURL("")
	a.Error(err)
	a.Equal(0, len(teachers3))
}
