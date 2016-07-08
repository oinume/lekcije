package middleware

import (
	"net/http"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/web"
	"goji.io"
	"golang.org/x/net/context"
)

func SetDbToContext(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		db, c, err := model.OpenAndSetTo(ctx)
		if err != nil {
			web.InternalServerError(w, err)
			return
		}
		defer db.Close()
		h.ServeHTTPC(c, w, r)
	}
	return goji.HandlerFunc(fn)
}
