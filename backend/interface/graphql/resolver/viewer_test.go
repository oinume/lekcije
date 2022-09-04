package resolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/di"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

func TestViewer(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	userUsecase := di.NewUserUsecase(db.DB())

	r := NewResolver(
		repos.FollowingTeacher(),
		repos.NotificationTimeSpan(),
		nil,
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)

	ctx := context.Background()
	user := modeltest.NewUser()
	repos.CreateUsers(ctx, t, user)
	userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
		uat.UserID = user.ID
	})
	repos.CreateUserAPITokens(ctx, t, userAPIToken)

	ctx = context_data.SetAPIToken(ctx, userAPIToken.Token)
	graphqlUser, err := r.Query().Viewer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("user = %+v\n", graphqlUser)
}
