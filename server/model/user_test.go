package model

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/server/util"
)

var _ = fmt.Print

func TestUserService_FindByGoogleID(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateRandomUser(t)
	userGoogle := helper.CreateUserGoogle(t, "1", user.ID)
	userActual, err := userService.FindByGoogleID(userGoogle.GoogleID)
	r.Nil(err)
	a.Equal(user.ID, userActual.ID)
	a.Equal(user.Email, userActual.Email)
}

func TestUserService_FindAllEmailVerifiedIsTrue(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	helper.TruncateAllTables(t)

	user := helper.CreateRandomUser(t)
	teacher := helper.CreateRandomTeacher(t)
	_ = helper.CreateFollowingTeacher(t, user.ID, teacher)

	users, err := userService.FindAllEmailVerifiedIsTrue(context.Background(), 10)
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

	user := helper.CreateRandomUser(t)
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
