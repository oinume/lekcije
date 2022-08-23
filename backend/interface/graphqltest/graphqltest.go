package graphqltest

import (
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/Khan/genqlient/graphql"

	"github.com/oinume/lekcije/backend/interface/graphql/generated"
	"github.com/oinume/lekcije/backend/interface/graphql/resolver"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/usecase"
)

func newResolver(repos *mysqltest.Repositories, userUsecase *usecase.User) *resolver.Resolver {
	return resolver.NewResolver(
		repos.FollowingTeacher(),
		repos.NotificationTimeSpan(),
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)
}

func newGraphQLServer(repos *mysqltest.Repositories, userUsecase *usecase.User) *handler.Server {
	resolver := newResolver(repos, userUsecase)
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
}
