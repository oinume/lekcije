package model

import (
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type User struct {
	ID                uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name              string
	Email             string
	RawEmail          string
	EmailVerified     bool
	PlanID            uint8
	FollowedTeacherAt mysql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
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

func (s *UserService) FindByUserAPIToken(userAPIToken string) (*User, error) {
	user := &User{}
	sql := `
	SELECT u.* FROM user AS u
	INNER JOIN user_api_token AS uat ON u.id = uat.user_id
	WHERE uat.token = ?
	`
	if result := s.db.Raw(sql, userAPIToken).Scan(user); result.Error != nil {
		if err := wrapNotFound(result, "User not found: userAPIToken=%v", userAPIToken); err != nil {
			return nil, err
		}
		return nil, errors.InternalWrapf(
			result.Error, "userAPIToken=%v", userAPIToken,
		)
	}
	return user, nil
}

// Returns an empty slice if no users found
func (s *UserService) FindAllEmailVerifiedIsTrue() ([]*User, error) {
	var users []*User
	sql := `SELECT * FROM user WHERE email_verified = 1 ORDER BY id`
	result := s.db.Raw(sql).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return nil, errors.InternalWrapf(result.Error, "Failed to find Users")
	}
	return users, nil
}

// Returns an empty slice if no users found
func (s *UserService) FindAllFollowedTeacherAtIsNull(createdAt time.Time) ([]*User, error) {
	var users []*User
	sql := `SELECT * FROM user WHERE followed_teacher_at IS NULL AND CAST(created_at AS DATE) = ? ORDER BY id`
	result := s.db.Raw(sql, createdAt.Format(dbDateFormat)).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return nil, errors.InternalWrapf(result.Error, "Failed to find Users")
	}
	return users, nil
}

func (s *UserService) Create(name, email string) (*User, error) {
	user := &User{
		Name:          name,
		Email:         email,
		RawEmail:      email,
		EmailVerified: true,
		PlanID:        DefaultPlanID,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, errors.InternalWrapf(result.Error, "")
	}
	return user, nil
}

func (s *UserService) CreateWithGoogle(name, email, googleID string) (*User, *UserGoogle, error) {
	user := &User{
		Name:          name,
		Email:         email,
		RawEmail:      email,
		EmailVerified: true, // TODO: set false after implement email verification
		PlanID:        DefaultPlanID,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, nil, errors.InternalWrapf(
			result.Error,
			"Failed to create User: email=%v, googleID=%v",
			email, googleID,
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
	result := s.db.Exec("UPDATE user SET email = ?, raw_email = ? WHERE id = ?", newEmail, newEmail, user.ID)
	if result.Error != nil {
		return errors.InternalWrapf(
			result.Error,
			"Failed to update email: id=%v, email=%v", user.ID, newEmail,
		)
	}
	return nil
}

func (s *UserService) UpdateFollowedTeacherAt(user *User) error {
	sql := "UPDATE user SET followed_teacher_at = NOW() WHERE id = ?"
	if err := s.db.Exec(sql, user.ID).Error; err != nil {
		return errors.InternalWrapf(err, "Failed to update followed_teacher_at: id=%v", user.ID)
	}
	return nil
}

func (s *UserService) FindLoggedInUser(token string) (*User, error) {
	user := &User{}
	sql := `
		SELECT * FROM user AS u
		INNER JOIN user_api_token AS uat ON u.id = uat.user_id
		WHERE uat.token = ?
		`
	result := s.db.Model(&User{}).Raw(strings.TrimSpace(sql), token).Scan(user)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, errors.NotFoundWrapf(result.Error, "Failed to find user: token=%s", token)
		}
		return nil, errors.InternalWrapf(result.Error, "find user: token=%s", token)
	}
	return user, nil
}
