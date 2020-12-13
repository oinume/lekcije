package e2e

import (
	"fmt"
	"net/http"

	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/util"
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

func getCookie(cookies []*http.Cookie, name string) (*http.Cookie, error) {
	for _, c := range cookies {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, http.ErrNoCookie
}
