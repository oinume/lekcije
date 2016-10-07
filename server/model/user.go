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
	ID            uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name          string
	Email         Email
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (*User) TableName() string {
	return "user"
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) TableName() string {
	return (&User{}).TableName()
}

func (s *UserService) FindByPk(id uint32) (*User, error) {
	user := &User{}
	if err := s.db.First(user, &User{ID: id}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByGoogleID(googleID string) (*User, error) {
	user := &User{}
	sql := `
	SELECT u.* FROM user AS u
	INNER JOIN user_google AS ug ON u.id = ug.user_id
	WHERE ug.google_id = ?
	LIMIT 1
	`
	if result := s.db.Raw(sql, googleID).Scan(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, errors.NotFoundWrapf(
				result.Error, "UserGoogle not found: googleID=%v", googleID,
			)
		} else {
			return nil, errors.InternalWrapf(
				result.Error, "googleID=%v", googleID,
			)
		}
	}
	return user, nil
}

func (s *UserService) FindOrCreate(name string, email Email) (*User, error) {
	user := User{
		Name:          name,
		Email:         email,
		EmailVerified: true, // TODO: set false after implement email verification
	}
	if err := s.db.FirstOrCreate(&user, User{Email: email}).Error; err != nil {
		return nil, errors.InternalWrapf(err, "Failed to get or create User: email=%v", email)
	}
	return &user, nil
}

func (s *UserService) Create(name, email string) (*User, error) {
	e, err := NewEmailFromRaw(email)
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:  name,
		Email: e,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, errors.InternalWrapf(result.Error, "")
	}
	return user, nil
}

func (s *UserService) CreateWithGoogle(name, email, googleID string) (*User, *UserGoogle, error) {
	e, err := NewEmailFromRaw(email)
	if err != nil {
		return nil, nil, err
	}

	user := &User{
		Name:  name,
		Email: e,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, nil, errors.InternalWrapf(
			result.Error, "Failed to create User: email=%v", email,
		)
	}

	userGoogle := &UserGoogle{
		GoogleID: googleID,
		UserID:   user.ID,
	}
	if result := s.db.Create(userGoogle); result.Error != nil {
		return nil, nil, errors.InternalWrapf(
			result.Error, "Failed to create UserGoogle: email=%v", email,
		)
	}

	return user, userGoogle, nil
}

func (s *UserService) UpdateEmail(user *User, newEmail string) error {
	email, err := NewEmailFromRaw(newEmail)
	if err != nil {
		return err
	}
	result := s.db.Exec("UPDATE user SET email = ? WHERE id = ?", email, user.ID)
	if result.Error != nil {
		return errors.InternalWrapf(
			result.Error,
			"Failed to update email: id=%v, email=%v", user.ID, email,
		)
	}
	return nil
}

func FindLoggedInUserAndSetToContext(token string, ctx context.Context) (*User, context.Context, error) {
	db := MustDB(ctx)
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
