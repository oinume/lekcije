package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

func main() {
	m := &followReminderMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type followReminderMain struct {
	outStream io.Writer
	errStream io.Writer
}

func (m *followReminderMain) run(args []string) error {
	flagSet := flag.NewFlagSet("follow_reminder", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		dryRun     = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
		targetDate = flagSet.String("target-date", "", "Specify registration date of users")
		logLevel   = flagSet.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	config.MustProcessDefault()
	startedAt := time.Now().UTC()
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel(*logLevel))
	appLogger.Info("follow_reminder started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		appLogger.Info("follow_reminder finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	var date time.Time
	if *targetDate == "" {
		date = time.Now().UTC().Add(-1 * 24 * time.Hour)
	} else {
		date, err = time.Parse("2006-01-02", *targetDate)
		if err != nil {
			return fmt.Errorf("failed to parse 'target-date': %v", *targetDate)
		}
	}
	users, err := model.NewUserService(db).FindAllFollowedTeacherAtIsNull(date)
	if err != nil {
		return err
	}

	sender := emailer.NewSendGridSender(http.DefaultClient, appLogger)
	templateText := getEmailTemplate()
	for _, user := range users {
		t := emailer.NewTemplate("follow_reminder", templateText)
		data := struct {
			To   string
			Name string
		}{
			To:   user.Email,
			Name: user.Name,
		}
		mail, err := emailer.NewEmailFromTemplate(t, data)
		if err != nil {
			return err
		}
		mail.SetCustomArg("email_type", model2.EmailTypeFollowReminder)
		mail.SetCustomArg("user_id", fmt.Sprint(user.ID))

		if !*dryRun {
			if err := sender.Send(ctx, mail); err != nil {
				return err
			}
		}
		appLogger.Info("followReminder", zap.Uint("userID", uint(user.ID)), zap.String("email", user.Email))
	}

	return nil
}

func getEmailTemplate() string {
	return strings.TrimSpace(`
From: lekcije@lekcije.com
To: {{ .To }}
Subject: lekcijeでお気に入りの講師をフォローしましょう
Body: text/html
<a href="https://www.lekcije.com/">lekcije</a>をご利用いただきありがとうございます。
lekcijeでお気入りのDMM英会話の講師をフォローしてみましょう。
フォローすると、講師の空きレッスンが登録された時にメールでお知らせが届きます。

フォローの仕方がわからない場合は下記をチェック！
<a href="https://lekcije.amebaownd.com/posts/2044879">PC</a>
<a href="https://lekcije.amebaownd.com/posts/1577091">Mobile</a>

今後ともlekcijeをよろしくお願いいたします。
	`)
}
