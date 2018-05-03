package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/controller/flash_message"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/event_logger"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/rs/cors"
	"go.uber.org/zap"
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
					err = fmt.Errorf("unknown error type: %v", errorType)
				}
				controller.InternalServerError(w, errors.NewInternalError(
					errors.WithError(err),
					errors.WithMessage("panic ocurred"),
				))
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
			trackingID := ""
			if v, err := context_data.GetTrackingID(r.Context()); err == nil {
				trackingID = v
			}

			// 180.76.15.26 - - [31/Jul/2016:13:18:07 +0000] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"
			logger.Access.Info(
				"access",
				zap.String("date", start.Format(time.RFC3339)),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", status),
				zap.Int("bytes", writerProxy.BytesWritten()),
				zap.String("remoteAddr", controller.GetRemoteAddress(r)),
				zap.String("userAgent", r.Header.Get("User-Agent")),
				zap.String("referer", r.Referer()),
				zap.Duration("elapsed", end.Sub(start)/time.Millisecond),
				zap.String("trackingID", trackingID),
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
		logger.App.Error("Failed to newrelic.NewApplication()", zap.Error(err))
		return h
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		tx := app.StartTransaction(r.URL.Path, w, r)
		defer tx.End()
		h.ServeHTTP(tx, r)
	}
	return http.HandlerFunc(fn)
}

func SetDBAndRedis(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.RequestURI == "/api/status" {
			h.ServeHTTP(w, r)
			return
		}
		if config.IsLocalEnv() {
			fmt.Printf("%s %s\n", r.Method, r.RequestURI)
		}

		db, err := model.OpenDB(
			config.DefaultVars.DBURL(),
			maxDBConnections,
			config.DefaultVars.DebugSQL,
		)
		if err != nil {
			controller.InternalServerError(w, err)
			return
		}
		defer db.Close()
		ctx = context_data.SetDB(ctx, db)

		redisClient, c, err := model.OpenRedisAndSetToContext(ctx, os.Getenv("REDIS_URL"))
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

func SetLoggedInUser(h http.Handler) http.Handler {
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

		userService := model.NewUserService(context_data.MustDB(ctx))
		user, err := userService.FindLoggedInUser(cookie.Value)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		c := context_data.SetLoggedInUser(ctx, user)
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func SetTrackingID(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ignoreURLs := []string{
			"/api/status",
			"/robots.txt",
			"/sitemap.xml",
		}
		for _, u := range ignoreURLs {
			if r.RequestURI == u {
				h.ServeHTTP(w, r)
				return
			}
		}

		cookie, err := r.Cookie(controller.TrackingIDCookieName)
		var trackingID string
		if err == nil {
			trackingID = cookie.Value
		} else {
			trackingID = uuid.New().String()
			domain := strings.Replace(r.Host, "www.", "", 1)
			domain = strings.Replace(domain, ":4000", "", 1) // TODO: local only
			c := &http.Cookie{
				Name:     controller.TrackingIDCookieName,
				Value:    trackingID,
				Path:     "/",
				Domain:   domain,
				Expires:  time.Now().UTC().Add(time.Hour * 24 * 365 * 2),
				HttpOnly: true,
			}
			http.SetCookie(w, c)
		}
		c := context_data.SetTrackingID(r.Context(), trackingID)
		r.Header.Set("Grpc-Metadata-Http-Tracking-Id", trackingID)
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func SetGRPCMetadata(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Grpc-Metadata-Http-User-Agent", r.UserAgent())
		r.Header.Set("Grpc-Metadata-Http-Referer", r.Referer())
		r.Header.Set("Grpc-Metadata-Http-Host", r.Host)
		r.Header.Set("Grpc-Metadata-Http-Url-Path", r.URL.Path)
		r.Header.Set("Grpc-Metadata-Http-Remote-Addr", getRemoteAddress(r))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func SetGAMeasurementEventValues(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		v := event_logger.NewGAMeasurementEventValuesFromRequest(r)
		c := event_logger.WithGAMeasurementEventValues(r.Context(), v)
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
			logger.App.Debug("Not logged in")
			http.Redirect(w, r, config.WebURL(), http.StatusFound)
			return
		}

		// TODO: Use context_data.MustLoggedInUser(ctx)
		userService := model.NewUserService(context_data.MustDB(ctx))
		user, err := userService.FindLoggedInUser(cookie.Value)
		if err != nil {
			if errors.IsNotFound(err) {
				logger.App.Debug("not logged in")
				http.Redirect(w, r, config.WebURL(), http.StatusFound)
				return
			}
			controller.InternalServerError(w, err)
			return
		}
		logger.App.Debug("Logged in user", zap.String("name", user.Name))
		c := context_data.SetLoggedInUser(ctx, user)
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func CORS(h http.Handler) http.Handler {
	origins := []string{}
	if strings.HasPrefix(config.StaticURL(), "http") {
		origins = append(origins, strings.TrimSuffix(config.StaticURL(), "/static"))
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

func getRemoteAddress(req *http.Request) string {
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		return (strings.Split(req.RemoteAddr, ":"))[0]
	}
	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
}
