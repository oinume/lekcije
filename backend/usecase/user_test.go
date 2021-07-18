package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
)

func Test_User_CreateWithGoogle(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	userUsecase := NewUser(repos.DB(), repos.User(), repos.UserGoogle())

	tests := map[string]struct {
		name     string
		email    string
		googleID string
	}{
		"normal": {
			name:     "oinume",
			email:    "oinume@gmail.com",
			googleID: "xyz",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			user, _, err := userUsecase.CreateWithGoogle(ctx, test.name, test.email, test.googleID)
			if err != nil {
				t.Fatal(err)
			}
			// TODO: validate user, userGoogle
			fmt.Printf("userID = %v\n", user.ID)
		})
	}
}
