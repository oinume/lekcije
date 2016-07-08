package middleware

import (
	"net/http"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/web"
	"goji.io"
	"golang.org/x/net/context"
)

func SetGormDbToContext(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		db, err := model.Open()
		if err != nil {
			web.InternalServerError(w, err)
			return
		}
		defer db.Close()
		c := context.WithValue(ctx, model.ContextKey, db)
		h.ServeHTTPC(c, w, r)
		//db2 := model.MustFromContext(c)
	}
	return goji.HandlerFunc(fn)
}
