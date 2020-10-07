package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/ga_measurement"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

var _ = fmt.Print

func panicHandler(h http.Handler) http.Handler {
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
				internalServerError(nil, w, errors.NewInternalError(
					errors.WithError(err),
					errors.WithMessage("panic occurred"),
				), 0)
				return
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func accessLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writerProxy := WrapWriter(w)
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
				logger.Info(
					"access",
					zap.String("date", start.Format(time.RFC3339)),
					zap.String("method", r.Method),
					zap.String("url", r.URL.String()),
					zap.Int("status", status),
					zap.Int("bytes", writerProxy.BytesWritten()),
					zap.String("remoteAddr", getRemoteAddress(r)),
					zap.String("userAgent", r.Header.Get("User-Agent")),
					zap.String("referer", r.Referer()),
					zap.Duration("elapsed", end.Sub(start)/time.Millisecond),
					zap.String("trackingID", trackingID),
				)
			}()
		}
		return http.HandlerFunc(fn)
	}
}

func setLoggedInUser(db *gorm.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if r.RequestURI == "/api/status" {
				h.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie(APITokenCookieName)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			userService := model.NewUserService(db)
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
}

func setTrackingID(h http.Handler) http.Handler {
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

		cookie, err := r.Cookie(TrackingIDCookieName)
		var trackingID string
		if err == nil {
			trackingID = cookie.Value
		} else {
			trackingID = uuid.New().String()
			domain := strings.Replace(r.Host, "www.", "", 1)
			domain = strings.Replace(domain, ":4000", "", 1) // TODO: local only
			c := &http.Cookie{
				Name:     TrackingIDCookieName,
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

func setGRPCMetadata(h http.Handler) http.Handler {
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

func setGAMeasurementEventValues(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c := ga_measurement.WithEventValues(
			r.Context(),
			ga_measurement.NewEventValuesFromRequest(r),
		)
		h.ServeHTTP(w, r.WithContext(c))
	}
	return http.HandlerFunc(fn)
}

func loginRequiredFilter(db *gorm.DB, appLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if !strings.HasPrefix(r.RequestURI, "/me") {
				h.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie(APITokenCookieName)
			if err != nil {
				appLogger.Debug("Not logged in")
				http.Redirect(w, r, config.WebURL(), http.StatusFound)
				return
			}

			// TODO: Use context_data.MustLoggedInUser(ctx)
			userService := model.NewUserService(db)
			user, err := userService.FindLoggedInUser(cookie.Value)
			if err != nil {
				if errors.IsNotFound(err) {
					appLogger.Debug("not logged in")
					http.Redirect(w, r, config.WebURL(), http.StatusFound)
					return
				}
				internalServerError(appLogger, w, err, 0)
				return
			}
			appLogger.Debug("Logged in user", zap.String("name", user.Name))
			c := context_data.SetLoggedInUser(ctx, user)
			h.ServeHTTP(w, r.WithContext(c))
		}
		return http.HandlerFunc(fn)
	}
}

func setCORS(h http.Handler) http.Handler {
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

func redirecter(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "lekcije.herokuapp.com" {
			http.Redirect(w, r, config.WebURL()+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
