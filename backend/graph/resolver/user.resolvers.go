package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oinume/lekcije/backend/graph/generated"
	"github.com/oinume/lekcije/backend/graph/model"
	"github.com/oinume/lekcije/backend/model2"
)

func (r *userResolver) FollowingTeachers(ctx context.Context, obj *model.User) ([]*model.FollowingTeacher, error) {
	userID, err := strconv.ParseUint(obj.ID, 10, 32)
	if err != nil {
		return nil, err
	}
	fts, err := r.followingTeacherRepo.FindByUserID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}
	followingTeachers := make([]*model.FollowingTeacher, len(fts))
	teacherIDs := make([]uint, len(fts))
	for i, ft := range fts {
		teacherIDs[i] = ft.TeacherID
	}
	teachers, err := r.teacherRepo.FindByIDs(ctx, teacherIDs)
	if err != nil {
		return nil, err
	}
	teachersMap := make(map[uint]*model2.Teacher, len(teachers))
	for _, t := range teachers {
		teachersMap[t.ID] = t
	}

	for i, ft := range fts {
		t, ok := teachersMap[ft.TeacherID]
		if !ok {
			continue
		}
		followingTeachers[i] = &model.FollowingTeacher{
			ID: fmt.Sprintf("%d_%d", ft.UserID, ft.TeacherID),
			Teacher: &model.Teacher{
				ID:   fmt.Sprint(t.ID),
				Name: t.Name,
			},
			CreatedAt: ft.CreatedAt.String(),
		}
	}

	return followingTeachers, nil
}

func (r *userResolver) NotificationTimeSpans(ctx context.Context, obj *model.User) ([]*model.NotificationTimeSpan, error) {
	panic(fmt.Errorf("not implemented"))
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
