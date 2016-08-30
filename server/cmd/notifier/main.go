package main

import (
	"flag"
	"os"
	"log"
	"context"
	"net/http"
	"fmt"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"time"
)

var (
	dryRun = flag.Bool("dry-run", false, "Dry run")
	logLevel = flag.String("log-level", "info", "Log level")
)

// TODO: move somewhere proper and make it be struct
var definedEnvs = map[string]string{
	"GOOGLE_CLIENT_ID":     "",
	"GOOGLE_CLIENT_SECRET": "",
	"DB_DSN":               "",
	"NODE_ENV":             "",
	"ENCRYPTION_KEY":       "",
}

var lessonFetcher *fetcher.TeacherLessonFetcher

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
	lessonFetcher = fetcher.NewTeacherLessonFetcher(http.DefaultClient, nil)
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
		os.Exit(1)
	}
}

func run() error {
	// Check env
	for key, _ := range definedEnvs {
		if value := os.Getenv(key); value != "" {
			definedEnvs[key] = value
		} else {
			log.Fatalf("Env '%v' must be defined.", key)
		}
	}

	db, _, err := model.OpenAndSetToContext(context.Background(), definedEnvs["DB_DSN"])
	if err != nil {
		return err
	}

	var users []*model.User
	// TODO: email_verified
	userSql := `SELECT * FROM user /* WHERE email_verified = 1 */`
	result := db.Raw(userSql).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return errors.InternalWrapf(result.Error, "")
	}

	for _, user := range users {
		teacherIds, err := model.FollowingTeacherService.FindTeacherIdsByUserId(user.Id)
		if err != nil {
			return err
		}
		for _, teacherId := range teacherIds {
			if err := fetchTeacherLessons(teacherId); err != nil {
				return err
			}
		}
	}

	return nil
}

func fetchTeacherLessons(teacherId uint32) error {
	teacher, lessons, err := lessonFetcher.Fetch(teacherId)
	if err != nil {
		return err
	}
	fmt.Printf("--- %s ---\n", teacher.Name)
	for _, lesson := range lessons {
		fmt.Printf("datetime = %v, status = %v\n", lesson.Datetime, lesson.Status)
	}

	// TODO: test
	_, err = model.LessonService.UpdateLessons(lessons)
	if err != nil {
		return err
	}
	return nil
}
