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
	user, err := UserService.Create("test", email)
	a.NoError(err)
	a.True(user.Id > 0)
	a.Equal(email, user.Email)
}

func TestUpdateEmail(t *testing.T) {
	a := assert.New(t)
	user := createTestUser()
	email := randomEmail()
	err := UserService.UpdateEmail(user, email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.NoError(err)

	actual, err := UserService.FindByPk(user.Id)
	a.NoError(err)
	a.NotEqual(user.Email, actual.Email)
}

func randomEmail() string {
	return util.RandomString(16) + "@example.com"
}

func createTestUser() *User {
	name := util.RandomString(16)
	email := name + "@example.com"
	user, err := UserService.Create(name, email)
	if err != nil {
		panic(err)
	}
	return user
}
