package web

import (
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
)

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(TemplateDir(), "index.html")))
	if err := t.Execute(w, nil); err != nil {
		panic(err)
	}
}

func PostMeFollowingTeachersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	teacherIdOrUrl := r.FormValue("teacherIdOrUrl")
	teacher, err := model.NewTeacherFromIdOrUrl(teacherIdOrUrl)
	if err != nil {
		InternalServerError(w, err)
		return
	}

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
