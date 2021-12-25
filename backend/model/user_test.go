package model

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/backend/randoms"
)

var _ = fmt.Print

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

func randomEmail() string {
	return randoms.MustNewString(16) + "@example.com"
}
