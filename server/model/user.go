package model

import (
	"strings"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"golang.org/x/net/context"
)

const (
	contextKeyLoggedInUser = "loggedInUser"
)

type User struct {
	Id        uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*User) TableName() string {
	return "user"
}

func (*UserApiToken) TableName() string {
	return "user_api_token"
}

type UserApiToken struct {
	Token     string `gorm:"primary_key;AUTO_INCREMENT"`
	UserId    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindLoggedInUserAndSetToContext(token string, ctx context.Context) (*User, context.Context, error) {
	db := MustDb(ctx)
	user := &User{}
	sql := `
		SELECT * FROM user AS u
		INNER JOIN user_api_token AS uat ON u.id = uat.user_id
		WHERE uat.token = ?
		`
	result := db.Model(&User{}).Raw(strings.TrimSpace(sql), token).Scan(user)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil, errors.NotFoundWrapf(result.Error, "Failed to find user: token=%s", token)
		}
		return nil, nil, errors.InternalWrapf(result.Error, "find user: token=%s", token)
	}
	c := context.WithValue(ctx, contextKeyLoggedInUser, user)
	return user, c, nil
}

// TODO: Move somewhere else model
func GetLoggedInUser(ctx context.Context) (*User, error) {
	value := ctx.Value(contextKeyLoggedInUser)
	if user, ok := value.(*User); ok {
		return user, nil
	}
	return nil, errors.NotFoundf("Logged in user not found in context")
}

func MustLoggedInUser(ctx context.Context) *User {
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}
