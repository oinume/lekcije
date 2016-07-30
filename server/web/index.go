package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/oinume/lekcije/server/errors"
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
	teacher, err := model.NewTeacherFromIdOrUrl(teacherIdOrUrl)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	// TODO: scraping teacher's lessons and insert teacher

	db := model.MustDb(ctx)
	now := time.Now()
	ft := model.FollowingTeacher{
		UserId:    user.Id,
		TeacherId: teacher.Id,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Save(ft).Error; err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
