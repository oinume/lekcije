package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/randoms"
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

func (s *UserAPITokenService) Create(userID uint32) (*UserAPIToken, error) {
	apiToken := randoms.MustNewString(64)
	userAPIToken := UserAPIToken{
		UserID: userID,
		Token:  apiToken,
	}
	if err := s.db.Create(&userAPIToken).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to create UserAPIToken"),
			errors.WithResource(errors.NewResource(userAPIToken.TableName(), "userID", userID)),
		)
	}
	return &userAPIToken, nil
}

func (s *UserAPITokenService) FindByPK(token string) (*UserAPIToken, error) {
	userAPIToken := &UserAPIToken{}
	if result := s.db.First(userAPIToken, &UserAPIToken{Token: token}); result.Error != nil {
		if err := wrapNotFound(result, userAPIToken.TableName(), "token", token); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(userAPIToken.TableName(), "token", token)),
		)
	}
	return userAPIToken, nil
}

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
