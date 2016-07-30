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
		fmt.Printf("%s %s\n", r.Method, r.RequestURI)

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

func SetLoggedInUserToContext(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/status" {
			h.ServeHTTPC(ctx, w, r)
			return
		}
		cookie, err := r.Cookie(web.ApiTokenCookieName)
		if err != nil {
			h.ServeHTTPC(ctx, w, r)
			return
		}

		user, c, err := model.FindLoggedInUserAndSetTo(cookie.Value, ctx)
		if err != nil {
			fmt.Printf("loggedInUser = %+v\n", user)
			h.ServeHTTPC(ctx, w, r)
			return
		}
		h.ServeHTTPC(c, w, r)
	}
	return goji.HandlerFunc(fn)
}
