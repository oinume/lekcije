package e2e

import (
	"fmt"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
)

func createUserAndLogin(name, email, googleID string) (*model.User, string, error) {
	userService := model.NewUserService(db)
	user, _, err := userService.CreateWithGoogle(name, email, googleID)
	if err != nil {
		return nil, "", err
	}
	apiToken, err := model.NewUserAPITokenService(db).Create(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, apiToken.Token, nil
}

func randomEmail(prefix string) string {
	return fmt.Sprintf("%s-%s@gmail.com", prefix, util.RandomString(8))
}
