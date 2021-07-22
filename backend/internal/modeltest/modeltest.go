package modeltest

import (
	"fmt"

	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/randoms"
)

func NewNotificationTimeSpan(setters ...func(nts *model2.NotificationTimeSpan)) *model2.NotificationTimeSpan {
	timeSpan := &model2.NotificationTimeSpan{}
	for _, setter := range setters {
		setter(timeSpan)
	}
	if timeSpan.UserID == 0 {
		timeSpan.UserID = uint(randoms.MustNewInt64(10000000))
	}
	if timeSpan.Number == 0 {
		timeSpan.Number = uint8(randoms.MustNewInt64(255))
	}
	if timeSpan.FromTime == "" {
		timeSpan.FromTime = ""
	}
	if timeSpan.ToTime == "" {
		timeSpan.ToTime = ""
	}
	return timeSpan
}

func NewUser(setters ...func(u *model2.User)) *model2.User {
	user := &model2.User{}
	for _, setter := range setters {
		setter(user)
	}
	if user.Name == "" {
		user.Name = "lekcije taro " + randoms.MustNewString(8)
	}
	if user.Email == "" {
		user.Email = fmt.Sprintf("lekcije-%s@example.com", randoms.MustNewString(8))
	}
	if user.PlanID == 0 {
		user.PlanID = uint8(model.DefaultMPlanID)
	}
	return user
}

func NewUserAPIToken(setters ...func(uat *model2.UserAPIToken)) *model2.UserAPIToken {
	userAPIToken := &model2.UserAPIToken{}
	for _, setter := range setters {
		setter(userAPIToken)
	}
	if userAPIToken.Token == "" {
		userAPIToken.Token = randoms.MustNewString(32)
	}
	if userAPIToken.UserID == 0 {
		userAPIToken.UserID = uint(randoms.MustNewInt64(10000000))
	}
	return userAPIToken
}

func NewUserGoogle(setters ...func(ug *model2.UserGoogle)) *model2.UserGoogle {
	userGoogle := &model2.UserGoogle{}
	for _, setter := range setters {
		setter(userGoogle)
	}
	if userGoogle.GoogleID == "" {
		userGoogle.GoogleID = randoms.MustNewString(32)
	}
	return userGoogle
}
