package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/oinume/lekcije/backend/interface/graphql/model"
)

// CreateFollowingTeacher is the resolver for the createFollowingTeacher field.
func (r *mutationResolver) CreateFollowingTeacher(ctx context.Context, input model.CreateFollowingTeacherInput) (*model.CreateFollowingTeacherPayload, error) {
	panic(fmt.Errorf("not implemented: CreateFollowingTeacher - createFollowingTeacher"))
}

// FollowingTeachers is the resolver for the followingTeachers field.
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
