package model

import (
	"fmt"
	"testing"

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

func TestUserService_FindAllEmailVerifiedIsTrue(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	helper.TruncateAllTables(helper.DB())

	user := helper.CreateRandomUser()
	teacher := helper.CreateRandomTeacher()
	_ = helper.CreateFollowingTeacher(user.ID, teacher)

	users, err := userService.FindAllEmailVerifiedIsTrue(10)
	r.NoError(err)
	a.Equal(1, len(users))
	a.Equal(user.ID, users[0].ID)
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
	r := require.New(t)
	email := randomEmail()
	user, err := userService.Create("test", email)
	r.NoError(err)
	a.True(user.ID > 0)
	a.Equal(email, user.Email)
	a.Equal(DefaultMPlanID, user.PlanID)
}

func TestUserService_UpdateEmail(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser()
	email := randomEmail()
	err := userService.UpdateEmail(user, email)
	r.NoError(err)

	actual, err := userService.FindByPK(user.ID)
	r.NoError(err)
	a.NotEqual(user.Email, actual.Email)
	a.Equal(email, actual.Email)
}

func randomEmail() string {
	return util.RandomString(16) + "@example.com"
}
