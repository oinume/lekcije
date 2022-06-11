package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/oinume/lekcije/backend/interface/graphql/model"
)

func (r *queryResolver) FollowingTeachers(ctx context.Context) ([]*model.FollowingTeacher, error) {
	//teachers, err := r.followingTeacherRepo.FindTeachersByUserID()
	return []*model.FollowingTeacher{
		{
			ID: "1",
			Teacher: &model.Teacher{
				ID:   "12345",
				Name: "oinume",
			},
			CreatedAt: "",
		},
	}, nil
}
