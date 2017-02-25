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

	userService := model.NewUserService(db)
	users, err := userService.FindAllFollowedTeacherAtIsNull()
	if err != nil {
		return err
	}

	sender := emailer.NewSendGridSender()
	for _, user := range users {
		t := emailer.NewTemplate("follow_reminder", user.Name)
		mail, err := emailer.NewEmailFromTemplate(t)
		if err != nil {
			return err
		}
		sender.Send(mail)
	}
	/*
			// 単数を送る場合
			t := email.Template()
			mail := email.NewEmailFromTemplate(t, templateValues)
			email.NewSender().Send(mail) // or mail.Send()

			// 複数送る場合
			t := email.Template()
			//mail := email.NewEmailFromTemplate(t, templateValues)
			sender := email.NewSender()
			mails := make([]a)
			for _, user := range users {
				mail := email.NewEmailFromTemplate(t, templateValues)
				mails = append(mails, mail)
			}
			email.NewSender().SendMulti(mails)

			return strings.TrimSpace(`
		From: hoge
		Subject: こんにちは
		Body:
		{{- range $teacherID := .TeacherIDs }}
		{{- $teacher := index $.Teachers $teacherID -}}
		--- {{ $teacher.Name }} ---
		  {{- $lessons := index $.LessonsPerTeacher $teacherID }}
		  {{- range $lesson := $lessons }}
		{{ $lesson.Datetime.Format "2006-01-02 15:04" }}
		  {{- end }}

		レッスンの予約はこちらから:
		<a href="http://eikaiwa.dmm.com/teacher/index/{{ $teacherID }}/">PC</a>
		<a href="http://eikaiwa.dmm.com/teacher/schedule/{{ $teacherID }}/">Mobile</a>

		{{ end }}
		空きレッスンの通知の解除は<a href="{{ .WebURL }}/me">こちら</a>
			`)
	*/

	return nil
}
