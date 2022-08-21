package resolver

import (
	"fmt"

	gqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model"
)

func toGraphQLUser(user *model.User) *gqlmodel.User {
	return &gqlmodel.User{
		ID:           fmt.Sprint(user.ID),
		Email:        user.Email,
		ShowTutorial: !user.IsFollowedTeacher(),
	}
}
