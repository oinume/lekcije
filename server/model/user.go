package model

import (
	"fmt"
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

func (s *UserService) FindByPK(id uint32) (*User, error) {
	user := &User{}
	if result := s.db.First(user, &User{ID: id}); result.Error != nil {
		if err := wrapNotFound(result, user.TableName(), "id", fmt.Sprint(id)); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(user.TableName(), "id", id)),
		)
	}
	if err := s.db.First(user, &User{ID: id}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByEmail(email string) (*User, error) {
	user := &User{}
	if result := s.db.First(user, &User{Email: email}); result.Error != nil {
		if err := wrapNotFound(result, user.TableName(), "email", email); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(user.TableName(), "email", email)),
		)
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
		if err := wrapNotFound(result, "user_google", "google_id", googleID); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource("user_google", "google_id", googleID)),
		)
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
		if err := wrapNotFound(result, user.TableName(), "userAPIToken", userAPIToken); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(user.TableName(), "userAPIToken", userAPIToken)),
		)
	}
	return user, nil
}

// Returns an empty slice if no users found
func (s *UserService) FindAllEmailVerifiedIsTrue(notificationInterval int) ([]*User, error) {
	var users []*User
	//sql := `
	//SELECT u.* FROM (SELECT DISTINCT(user_id) FROM following_teacher) AS ft
	//INNER JOIN user AS u ON ft.user_id = u.id
	//INNER JOIN m_plan AS mp ON u.plan_id = mp.id
	//WHERE
	//  u.email_verified = 1
	//  AND mp.notification_interval = ?
	//`
	sql := `
SELECT u.* FROM users AS u
INNER JOIN m_plan AS mp ON u.plan_id = mp.id
WHERE
  u.email_verified = 1
  AND mp.notification_interval = ?
`
	result := s.db.Raw(strings.TrimSpace(sql), notificationInterval).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to find Users"),
		)
	}
	return users, nil
}

// Returns an empty slice if no users found
func (s *UserService) FindAllFollowedTeacherAtIsNull(createdAt time.Time) ([]*User, error) {
	var users []*User
	sql := `SELECT * FROM user WHERE followed_teacher_at IS NULL AND CAST(created_at AS DATE) = ? ORDER BY id`
	result := s.db.Raw(sql, createdAt.Format(dbDateFormat)).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to find users"),
		)
	}
	return users, nil
}

func (s *UserService) Create(name, email string) (*User, error) {
	user := &User{
		Name:          name,
		Email:         email,
		EmailVerified: true,
		PlanID:        DefaultMPlanID,
	}
	if result := s.db.Create(user); result.Error != nil {
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to Create user"),
			errors.WithResource(errors.NewResource("user", "email", email)),
		)
	}
	return user, nil
}

func (s *UserService) CreateWithGoogle(name, email, googleID string) (*User, *UserGoogle, error) {
	user, err := s.FindByEmail(email)
	if e, ok := err.(*errors.AnnotatedError); ok && e.IsNotFound() {
		user = &User{
			Name:          name,
			Email:         email,
			EmailVerified: true,
			PlanID:        DefaultMPlanID,
		}
		if result := s.db.Create(user); result.Error != nil {
			return nil, nil, errors.NewInternalError(
				errors.WithError(result.Error),
				errors.WithMessage("Failed to create User"),
				errors.WithResource(errors.NewResourceWithEntries(
					"user", []errors.ResourceEntry{
						{"email", email}, {"googleID", googleID},
					},
				)),
			)
		}
	}
	// Do nothing if the user exists.

	userGoogleService := NewUserGoogleService(s.db)
	userGoogle, err := userGoogleService.FindByUserID(user.ID)
	if e, ok := err.(*errors.AnnotatedError); ok && e.IsNotFound() {
		userGoogle = &UserGoogle{
			GoogleID: googleID,
			UserID:   user.ID,
		}
		if result := s.db.Create(userGoogle); result.Error != nil {
			return nil, nil, errors.NewInternalError(
				errors.WithError(result.Error),
				errors.WithMessage("Failed to create UserGoogle"),
				errors.WithResource(errors.NewResource("user_google", "googleID", googleID)),
			)
		}
	}
	// Do nothing if the user google exists.

	return user, userGoogle, nil
}

func (s *UserService) UpdateEmail(user *User, newEmail string) error {
	result := s.db.Exec("UPDATE user SET email = ? WHERE id = ?", newEmail, user.ID)
	if result.Error != nil {
		return errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to update user email"),
			errors.WithResource(errors.NewResourceWithEntries(
				user.TableName(), []errors.ResourceEntry{
					{"id", user.ID}, {"email", newEmail},
				},
			)),
		)
	}
	return nil
}

func (s *UserService) UpdateFollowedTeacherAt(user *User) error {
	sql := "UPDATE user SET followed_teacher_at = NOW() WHERE id = ?"
	if err := s.db.Exec(sql, user.ID).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to update followed_teacher_at"),
			errors.WithResource(errors.NewResource(user.TableName(), "id", user.ID)),
		)
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
		if err := wrapNotFound(result, user.TableName(), "token", token); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(user.TableName(), "token", token)),
		)
	}
	return user, nil
}
