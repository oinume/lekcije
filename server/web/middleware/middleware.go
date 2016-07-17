package middleware

import (
	"net/http"

	"fmt"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/web"
	"goji.io"
	"golang.org/x/net/context"
)

var _ = fmt.Print

func SetDbToContext(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/status" {
			h.ServeHTTPC(ctx, w, r)
			return
		}

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
