package context_data

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
)

type dbKey struct{}
type loggedInUserKey struct{}
type trackingIDKey struct{}

func SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

func GetDB(ctx context.Context) (*gorm.DB, error) {
	value := ctx.Value(dbKey{})
	if db, ok := value.(*gorm.DB); ok {
		return db, nil
	}
	return nil, errors.NotFoundf("Failed to get *gorm.DB from context")
}

func MustDB(ctx context.Context) *gorm.DB {
	db, err := GetDB(ctx)
	if err != nil {
		panic(err)
	}
	return db
}

func GetLoggedInUser(ctx context.Context) (*model.User, error) {
	value := ctx.Value(loggedInUserKey{})
	if user, ok := value.(*model.User); ok {
		return user, nil
	}
	return nil, errors.NotFoundf("Failed to get loggedInUser from context")
}

func MustLoggedInUser(ctx context.Context) *model.User {
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

func SetLoggedInUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, loggedInUserKey{}, user)
}

func SetTrackingID(ctx context.Context, trackingID string) context.Context {
	return context.WithValue(ctx, trackingIDKey{}, trackingID)
}

func MustTrackingID(ctx context.Context) string {
	value := ctx.Value(trackingIDKey{})
	if trackingID, ok := value.(string); ok {
		return trackingID
	} else {
		panic(fmt.Sprintf("Failed to get trackingID from context"))
	}
}
