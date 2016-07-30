package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/oinume/lekcije/server/model"
	"golang.org/x/net/context"
)

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(TemplateDir(), "index.html")))
	if err := t.Execute(w, nil); err != nil {
		panic(err)
	}
}

func PostFollowingTeachersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	teacherIdOrUrl := r.FormValue("teacherIdOrUrl")
	fmt.Printf("teacherIdOrUrl = %s\n", teacherIdOrUrl)
	fmt.Fprintf(w, teacherIdOrUrl)
	teacher, err := model.NewTeacherFromIdOrUrl(teacherIdOrUrl)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	db := model.MustDbFromContext(ctx)
	ft := model.FollowingTeacher{
		UserId:    1,
		TeacherId: teacher.Id,
	}
	if err := db.Save(ft).Error; err != nil {
		InternalServerError(w, err)
		return
	}
}
