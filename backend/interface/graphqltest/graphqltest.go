package graphqltest

import (
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/Khan/genqlient/graphql"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/interface/graphql/generated"
	"github.com/oinume/lekcije/backend/interface/graphql/resolver"
	interfacehttp "github.com/oinume/lekcije/backend/interface/http"
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

func setAuthorizationContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		auth, err := interfacehttp.ParseAuthorizationHeader(r.Header.Get("authorization"))
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context_data.SetAPIToken(r.Context(), strings.TrimSpace(auth)))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
