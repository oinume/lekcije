package usecase

import (
	"context"
	"fmt"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type User struct {
	dbRepo         repository.DB
	userRepo       repository.User
	userGoogleRepo repository.UserGoogle
}

func NewUser(dbRepo repository.DB, userRepo repository.User, userGoogleRepo repository.UserGoogle) *User {
	return &User{
		dbRepo:         dbRepo,
		userRepo:       userRepo,
		userGoogleRepo: userGoogleRepo,
	}
}

func (u *User) CreateWithGoogle(ctx context.Context, name, email, googleID string) (*model2.User, *model2.UserGoogle, error) {
	var (
		retUser       *model2.User
		retUserGoogle *model2.UserGoogle
	)
	if err := u.dbRepo.Transaction(ctx, func(exec repository.Executor) error {
		user, err := u.userRepo.FindByEmailWithExec(ctx, exec, email)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}

			// Create User
			user = &model2.User{
				Name:          name,
				Email:         email,
				EmailVerified: 1,
				PlanID:        model.DefaultMPlanID, // TODO: define in model2
			}
			if err := u.userRepo.CreateWithExec(ctx, exec, user); err != nil {
				return err
			}
			retUser = user

			// Create UserGoogle
			userGoogle := &model2.UserGoogle{
				UserID:   user.ID,
				GoogleID: googleID,
			}
			if err := u.userGoogleRepo.CreateWithExec(ctx, exec, userGoogle); err != nil {
				return err
			}
			retUserGoogle = userGoogle
		} else {
			retUser = user
		}

		existingUserGoogle, err := u.userGoogleRepo.FindByPKWithExec(ctx, exec, googleID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}

			if user == nil {
				return fmt.Errorf("CreateWithGoogle: must not reach here")
			} else {
				// User exists, it means the user is created with another google id somehow
				anotherUserGoogle, err := u.userGoogleRepo.FindByUserIDWithExec(ctx, exec, user.ID)
				if err != nil {
					if !errors.IsNotFound(err) {
						return err
					}
					userGoogle := &model2.UserGoogle{
						UserID:   user.ID,
						GoogleID: googleID,
					}
					if err := u.userGoogleRepo.CreateWithExec(ctx, exec, userGoogle); err != nil {
						return err
					}
				}

				// Delete existing UserGoogle first then create a new one
				if err := u.userGoogleRepo.DeleteByPKWithExec(ctx, exec, anotherUserGoogle.GoogleID); err != nil {
					return err
				}
				userGoogle := &model2.UserGoogle{
					UserID:   user.ID,
					GoogleID: googleID,
				}
				if err := u.userGoogleRepo.CreateWithExec(ctx, exec, userGoogle); err != nil {
					return err
				}
				retUserGoogle = userGoogle
			}
		} else {
			retUserGoogle = existingUserGoogle
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return retUser, retUserGoogle, nil
}

func (u *User) FindByGoogleID(ctx context.Context, googleID string) (*model2.User, error) {
	return u.userRepo.FindByGoogleID(ctx, googleID)
}

func (u *User) UpdateEmail(ctx context.Context, id uint, email string) error {
	return u.userRepo.UpdateEmail(ctx, id, email)
}

func (u *User) IsDuplicateEmail(ctx context.Context, email string) (bool, error) {
	_, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
