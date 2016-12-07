package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/newrelic/go-agent"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/flash_message"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/rs/cors"
	"github.com/uber-go/zap"
)

var _ = fmt.Print

const maxDBConnections = 5

func PanicHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch errorType := r.(type) {
				case string:
					err = fmt.Errorf(errorType)
				case error:
					err = errorType
				default:
					err = fmt.Errorf("Unknown error type: %v", errorType)
				}
				controller.InternalServerError(w, errors.InternalWrapf(err, "panic ocurred"))
				return
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func AccessLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writerProxy := controller.WrapWriter(w)
		h.ServeHTTP(writerProxy, r)
		func() {
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
		}()
	}
	return http.HandlerFunc(fn)
}

func NewRelic(h http.Handler) http.Handler {
	key := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if key == "" {
		return h
	}

	c := newrelic.NewConfig("lekcije", key)
	app, err := newrelic.NewApplication(c)
	if err != nil {
		logger.AppLogger.Error("Failed to newrelic.NewApplication()", zap.Error(err))
		return h
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		tx := app.StartTransaction(r.URL.Path, w, r)
		defer tx.End()
		h.ServeHTTP(tx, r)
	}
	return http.HandlerFunc(fn)
}

func SetDBAndRedisToContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.RequestURI == "/api/status" {
			h.ServeHTTP(w, r)
			return
		}
		if config.IsLocalEnv() {
			fmt.Printf("%s %s\n", r.Method, r.RequestURI)
		}

		db, c, err := model.OpenDBAndSetToContext(
			ctx, os.Getenv("DB_URL"), maxDBConnections, !config.IsProductionEnv(),
		)
		if err != nil {
			controller.InternalServerError(w, err)
			return
		}
		defer db.Close()

		redisClient, c, err := model.OpenRedisAndSetToContext(c, os.Getenv("REDIS_URL"))
		if err != nil {
			controller.InternalServerError(w, err)
			return
		}
		defer redisClient.Close()

		_, c = flash_message.NewStoreRedisAndSetToContext(c, redisClient)

		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func SetLoggedInUserToContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.RequestURI == "/api/status" {
			h.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie(controller.APITokenCookieName)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		user, c, err := model.FindLoggedInUserAndSetToContext(cookie.Value, ctx)
		if err != nil {
			fmt.Printf("loggedInUser = %+v\n", user)
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func LoginRequiredFilter(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !strings.HasPrefix(r.RequestURI, "/me") {
			h.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie(controller.APITokenCookieName)
		if err != nil {
			logger.AppLogger.Debug("Not logged in")
			http.Redirect(w, r, config.WebURL(), http.StatusFound)
			return
		}

		user, c, err := model.FindLoggedInUserAndSetToContext(cookie.Value, ctx)
		if err != nil {
			switch err.(type) {
			case *errors.NotFound:
				logger.AppLogger.Debug("not logged in")
				http.Redirect(w, r, config.WebURL(), http.StatusFound)
				return
			default:
				controller.InternalServerError(w, err)
				return
			}
		}
		logger.AppLogger.Debug("Logged in user", zap.Object("user", user))
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func CORS(h http.Handler) http.Handler {
	origins := []string{}
	if strings.HasPrefix(config.StaticURL(), "http") {
		origins = append(origins, strings.TrimRight(config.StaticURL(), "/static"))
	}
	c := cors.New(cors.Options{
		AllowedOrigins: origins,
		//Debug:          true,
	})
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.HandlerFunc(w, r)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Redirecter(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "lekcije.herokuapp.com" {
			http.Redirect(w, r, config.WebURL()+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
