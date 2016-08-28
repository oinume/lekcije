package model

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
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

type UserServiceType struct {
	db *gorm.DB
}

var UserService UserServiceType

func (s *UserServiceType) TableName() string {
	return (&User{}).TableName()
}

func (s *UserServiceType) CreateUser(name, email string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, errors.InternalWrapf(result.Error, "")
	}
	return user, nil
}

func (s *UserServiceType) UpdateEmail(user *User, newEmail string) error {
	email, err := NewEmailFromRaw(newEmail)
	if err != nil {
		return err
	}
	result := s.db.Exec("UPDATE user SET email = ? WHERE id = ?", email, user.Id)
	if result.Error != nil {
		return errors.InternalWrapf(
			result.Error,
			"Failed to update email: id=%v, email=%v", user.Id, email,
		)
	}
	return nil
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
