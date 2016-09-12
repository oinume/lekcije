package model

import (
	"fmt"
	"testing"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestCreateUser(t *testing.T) {
	a := assert.New(t)
	email := randomEmail()
	user, err := userService.Create("test", email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.NoError(err)
	a.True(user.Id > 0)
	a.Equal(email, user.Email.Raw())
}

func TestUpdateEmail(t *testing.T) {
	a := assert.New(t)
	user := createTestUser()
	email := randomEmail()
	err := userService.UpdateEmail(user, email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.NoError(err)

	actual, err := userService.FindByPk(user.Id)
	a.NoError(err)
	a.NotEqual(user.Email.Raw(), actual.Email.Raw())
	a.Equal(email, actual.Email.Raw())
}

func randomEmail() string {
	return util.RandomString(16) + "@example.com"
}

func createTestUser() *User {
	name := util.RandomString(16)
	email := name + "@example.com"
	user, err := userService.Create(name, email)
	if err != nil {
		panic(err)
	}
	return user
}
