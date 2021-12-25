package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
)

const UserAPITokenExpiration = time.Hour * 24 * 30

type UserAPIToken struct {
	Token     string `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*UserAPIToken) TableName() string {
	return "user_api_token"
}

type UserAPITokenService struct {
	db *gorm.DB
}

func NewUserAPITokenService(db *gorm.DB) *UserAPITokenService {
	return &UserAPITokenService{db: db}
}

// DeleteByUserIDAndToken is used from getMeLogout on me.go
func (s *UserAPITokenService) DeleteByUserIDAndToken(userID uint32, token string) error {
	result := s.db.Where("user_id = ? AND token = ?", userID, token).Delete(&UserAPIToken{})
	if result.Error != nil {
		return errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to DeleteByUserIDAndToken"),
			errors.WithResource(
				errors.NewResourceWithEntries(
					(&UserAPIToken{}).TableName(), []errors.ResourceEntry{
						{Key: "userID", Value: userID},
						{Key: "token", Value: token},
					},
				),
			),
		)
	}
	return nil
}
