package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/rs/cors"
	"github.com/uber-go/zap"
	"goji.io"
	"golang.org/x/net/context"
)

var _ = fmt.Print

func AccessLogger(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writerProxy := controller.WrapWriter(w)
		// TODO: handle panic
		h.ServeHTTPC(ctx, writerProxy, r)

		end := time.Now()
		status := writerProxy.Status()
		if status == 0 {
			status = http.StatusOK
		}
		remoteAddr := r.RemoteAddr
		if remoteAddr != "" {
			remoteAddr = (strings.Split(remoteAddr, ":"))[0]
		}

		// 180.76.15.26 - - [31/Jul/2016:13:18:07 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"
		logger.AccessLogger.Info(
			"",
			zap.String("date", start.Format(time.RFC3339)),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", status),
			zap.Int("bytes", writerProxy.BytesWritten()),
			zap.String("remoteAddr", remoteAddr),
			zap.String("userAgent", r.Header.Get("User-Agent")),
			zap.String("referer", r.Referer()),
			zap.Duration("elapsed", end.Sub(start)/time.Millisecond),
		)
	}
	return goji.HandlerFunc(fn)
}

func SetDbToContext(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/status" {
			h.ServeHTTPC(ctx, w, r)
			return
		}
		fmt.Printf("%s %s\n", r.Method, r.RequestURI)

		db, c, err := model.OpenAndSetToContext(ctx, os.Getenv("DB_DSN"))
		if err != nil {
			controller.InternalServerError(w, err)
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
		cookie, err := r.Cookie(controller.ApiTokenCookieName)
		if err != nil {
			h.ServeHTTPC(ctx, w, r)
			return
		}

		user, c, err := model.FindLoggedInUserAndSetToContext(cookie.Value, ctx)
		if err != nil {
			fmt.Printf("loggedInUser = %+v\n", user)
			h.ServeHTTPC(ctx, w, r)
			return
		}
		h.ServeHTTPC(c, w, r)
	}
	return goji.HandlerFunc(fn)
}

func LoginRequiredFilter(h goji.Handler) goji.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RequestURI, "/me") {
			h.ServeHTTPC(ctx, w, r)
			return
		}
		cookie, err := r.Cookie(controller.ApiTokenCookieName)
		if err != nil {
			h.ServeHTTPC(ctx, w, r)
			return
		}

		user, c, err := model.FindLoggedInUserAndSetToContext(cookie.Value, ctx)
		if err != nil {
			fmt.Printf("loggedInUser = %+v\n", user)
			h.ServeHTTPC(ctx, w, r)
			return
		}
		h.ServeHTTPC(c, w, r)
	}
	return goji.HandlerFunc(fn)
}

func CORS(h goji.Handler) goji.Handler {
	origins := []string{}
	if strings.HasPrefix(config.StaticURL(), "http") {
		origins = append(origins, strings.TrimRight(config.StaticURL(), "/static"))
	}
	c := cors.New(cors.Options{
		AllowedOrigins: origins,
		//Debug:          true,
	})
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		c.HandlerFunc(w, r)
		h.ServeHTTPC(ctx, w, r)
	}
	return goji.HandlerFunc(fn)
}
