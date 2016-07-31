package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
)

var _ = fmt.Print

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(TemplateDir(), "index.html")))
	type Data struct {
		Teachers []*model.Teacher
	}
	data := &Data{}

	user, err := model.GetLoggedInUser(ctx)
	fmt.Printf("user = %+v\n", user)
	if err == nil {
		teachers, err := model.FollowingTeacherRepo.FindTeachersByUserId(user.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		data.Teachers = teachers
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func PostMeFollowingTeachersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	teacherIdOrUrl := r.FormValue("teacherIdOrUrl")
	if teacherIdOrUrl == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, err := model.NewTeacherFromIdOrUrl(teacherIdOrUrl)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	fetcher := fetcher.NewTeacherLessonFetcher(http.DefaultClient, logger.AppLogger)
	teacher, _, err := fetcher.Fetch(t.Id)

	now := time.Now()
	teacher.CreatedAt = now
	teacher.UpdatedAt = now
	db := model.MustDb(ctx)
	if err := db.FirstOrCreate(teacher).Error; err != nil {
		e := errors.InternalWrapf(err, "Failed to create Teacher: teacherId=%d", teacher.Id)
		InternalServerError(w, e)
	}

	ft := model.FollowingTeacher{
		UserId:    user.Id,
		TeacherId: teacher.Id,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.FirstOrCreate(ft).Error; err != nil {
		e := errors.InternalWrapf(
			err,
			"Failed to create FollowingTeacher: userId=%d, teacherId=%d",
			user.Id, teacher.Id,
		)
		InternalServerError(w, e)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
