package model

import (
	"fmt"
	"testing"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = fmt.Print

func TestUserService_FindByGoogleID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	helper := NewTestHelper(t)

	user := helper.CreateRandomUser()
	userGoogle := helper.CreateUserGoogle("1", user.ID)
	userActual, err := userService.FindByGoogleID(userGoogle.GoogleID)
	r.Nil(err)
	a.Equal(user.ID, userActual.ID)
	a.Equal(user.Email, userActual.Email)
}

func TestCreateUser(t *testing.T) {
	a := assert.New(t)
	email := randomEmail()
	user, err := userService.Create("test", email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.Nil(err)
	a.True(user.ID > 0)
	a.Equal(email, user.Email)
	a.Equal(email, user.RawEmail)
	a.Equal(DefaultPlanID, user.PlanID)
}

func TestUpdateEmail(t *testing.T) {
	a := assert.New(t)
	helper := NewTestHelper(t)

	user := helper.CreateRandomUser()
	email := randomEmail()
	err := userService.UpdateEmail(user, email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.Nil(err)

	actual, err := userService.FindByPK(user.ID)
	a.Nil(err)
	a.NotEqual(user.Email, actual.Email)
	a.Equal(email, actual.Email)
}

func randomEmail() string {
	return util.RandomString(16) + "@example.com"
}
