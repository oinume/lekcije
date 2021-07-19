package modeltest

import (
	"fmt"

	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/util"
)

func NewUser(setters ...func(u *model2.User)) *model2.User {
	user := &model2.User{}
	for _, setter := range setters {
		setter(user)
	}
	if user.Name == "" {
		user.Name = "lekcije taro " + util.RandomString(8)
	}
	if user.Email == "" {
		user.Email = fmt.Sprintf("lekcije-%s@example.com", util.RandomString(8))
	}
	if user.PlanID == 0 {
		user.PlanID = uint8(model.DefaultMPlanID)
	}
	return user
}

func NewUserGoogle(setters ...func(ug *model2.UserGoogle)) *model2.UserGoogle {
	userGoogle := &model2.UserGoogle{}
	for _, setter := range setters {
		setter(userGoogle)
	}
	if userGoogle.GoogleID == "" {
		userGoogle.GoogleID = util.RandomString(32)
	}
	return userGoogle
}
