package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
)

type User struct {
	ID                 uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name               string
	Email              string
	EmailVerified      bool
	PlanID             uint8
	FollowedTeacherAt  sql.NullTime
	OpenNotificationAt sql.NullTime
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (*User) TableName() string {
	return "user"
}

func (u *User) IsFollowedTeacher() bool {
	return u.FollowedTeacherAt.Valid && !u.FollowedTeacherAt.Time.IsZero()
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

// FindAllFollowedTeacherAtIsNull returns an empty slice if no users found
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

func (s *UserService) UpdateOpenNotificationAt(userID uint32, t time.Time) error {
	sql := "UPDATE user SET open_notification_at = ? WHERE id = ?"
	if err := s.db.Exec(sql, t.Format(dbDatetimeFormat), userID).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to update open_notification_at"),
			errors.WithResource(errors.NewResource((&User{}).TableName(), "id", userID)),
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
