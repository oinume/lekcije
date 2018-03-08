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

	user := helper.CreateRandomUser()
	userGoogle := helper.CreateUserGoogle("1", user.ID)
	userActual, err := userService.FindByGoogleID(userGoogle.GoogleID)
	r.Nil(err)
	a.Equal(user.ID, userActual.ID)
	a.Equal(user.Email, userActual.Email)
}

func TestUserService_CreateWithGoogle(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	email := randomEmail()
	googleID := util.RandomString(16)
	user, userGoogle, err := userService.CreateWithGoogle(googleID, email, googleID)
	r.NoError(err)
	a.Equal(email, user.Email)
	a.Equal(user.ID, userGoogle.UserID)
	a.Equal(googleID, userGoogle.GoogleID)

	googleID2 := util.RandomString(16)
	user2, _, err := userService.CreateWithGoogle(googleID2, email, googleID2)
	r.NoError(err)
	a.Equal(user.ID, user2.ID)
	a.Equal(user.Email, user2.Email)
}

func TestUserService_Create(t *testing.T) {
	a := assert.New(t)
	email := randomEmail()
	user, err := userService.Create("test", email)
	if e, ok := err.(*errors.Internal); ok {
		fmt.Printf("%+v\n", e.StackTrace())
	}
	a.Nil(err)
	a.True(user.ID > 0)
	a.Equal(email, user.Email)
	a.Equal(DefaultMPlanID, user.PlanID)
}

func TestUserService_UpdateEmail(t *testing.T) {
	a := assert.New(t)

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
