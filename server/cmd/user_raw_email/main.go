package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/model"
)

var (
	run = flag.Bool("run", false, "Run flag")
)

func main() {
	flag.Parse()
	bootstrap.CheckCLIEnvVars()

	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL, 1, !config.IsProductionEnv())
	if err != nil {
		log.Fatalf("OpenDB: err=%v", err)
	}
	defer db.Close()

	userService := model.NewUserService(db)
	users, err := userService.FindAllEmailVerifiedIsTrue()
	if err != nil {
		log.Fatalf("FindAllEmailVerifiedIsTrue: err = %v", err)
	}

	for _, user := range users {
		rawEmail := user.Email.Raw()
		fmt.Printf("userID=%v, rawEmail=%v\n", user.ID, rawEmail)
		if *run {
			if err := userService.UpdateEmail(user, rawEmail); err != nil {
				log.Fatalf("err = %v", err)
			}
		}
	}
}
