package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/uber-go/zap"
	"net/http"
	"strings"
)

var (
	dryRun   = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	logLevel = flag.String("log-level", "info", "Log level")
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}

func run() error {
	bootstrap.CheckCLIEnvVars()
	startedAt := time.Now().UTC()
	if *logLevel != "" {
		logger.App.SetLevel(logger.NewLevel(*logLevel))
	}
	logger.App.Info("follow_reminder started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("follow_reminder finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL, 1, !config.IsProductionEnv())
	if err != nil {
		return err
	}
	defer db.Close()

	users, err := model.NewUserService(db).FindAllFollowedTeacherAtIsNull()
	if err != nil {
		return err
	}

	sender := emailer.NewSendGridSender(http.DefaultClient)
	templateText := getEmailTemplate()
	for _, user := range users {
		t := emailer.NewTemplate("follow_reminder", templateText)
		data := struct {
			To string
			Name string
		}{
			To: user.Email.Raw(),
			Name: user.Name,
		}
		mail, err := emailer.NewEmailFromTemplate(t, data)
		if err != nil {
			return err
		}

		if !*dryRun {
			if err := sender.Send(mail); err != nil {
				return err
			}
		}
	}

	return nil
}

func getEmailTemplate() string {
	return strings.TrimSpace(`
From: lekcije@lekcije.com
To: {{ .To }}
Subject: lekcijeでお気に入りの講師をフォローしましょう
Body: text/html
{{ .Name }} 様

<a href="https://www.lekcije.com/">lekcije</a>をご利用いただきありがとうございます。
lekcijeでお気入りのDMM英会話の講師をフォローしてみましょう。
フォローすると、講師の空きレッスンが登録された時にメールでお知らせが届きます。

フォローの仕方がわからない場合は下記をチェック！
<a href="https://lekcije.amebaownd.com/posts/2044879">PC</a>
<a href="https://lekcije.amebaownd.com/posts/1577091">Mobile</a>

今後ともlekcijeをよろしくお願いいたします。
	`)
}
