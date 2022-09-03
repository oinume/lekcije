package graphqltest

import (
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/Khan/genqlient/graphql"

	"github.com/oinume/lekcije/backend/interface/graphql/generated"
	"github.com/oinume/lekcije/backend/interface/graphql/resolver"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/usecase"
)

func newGraphQLServer(
	repos *mysqltest.Repositories,
	notificationTimeSpanUsecase *usecase.NotificationTimeSpan,
	userUsecase *usecase.User,
) *handler.Server {
	resolver := resolver.NewResolver(
		repos.FollowingTeacher(),
		repos.NotificationTimeSpan(),
		notificationTimeSpanUsecase,
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
}
