package controller

import (
	"html/template"
	"net/http"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"golang.org/x/net/context"
)

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
		return
	}

	ft := &model.FollowingTeacher{
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

func PostMeFollowingTeachersDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	teacherIds := r.Form["teacherIds"]
	if len(teacherIds) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	_, err := model.FollowingTeacherRepo.DeleteTeachersByUserIdAndTeacherIds(
		user.Id,
		util.StringToUint32Slice(teacherIds...),
	)
	if err != nil {
		e := errors.InternalWrapf(err, "Failed to delete Teachers: teacherIds=%v", teacherIds)
		InternalServerError(w, e)
		return
	}

	// TODO: stash
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetMeSetting(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	t := template.Must(template.ParseFiles(
		TemplatePath("_base.html"),
		TemplatePath("me/setting.html")),
	)
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(),
		Email:              user.Email,
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func PostMeSettingUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//user := model.MustLoggedInUser(ctx)
	http.Redirect(w, r, "/me/setting", http.StatusFound)
}
