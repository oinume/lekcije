package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"golang.org/x/net/context"
)

var _ = fmt.Print

type commonTemplateData struct {
	StaticUrl string
}

func getCommonTemplateData() commonTemplateData {
	return commonTemplateData{
		StaticUrl: config.StaticUrl(),
	}
}

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := model.GetLoggedInUser(ctx)
	if err == nil {
		indexLogin(ctx, w, r, user)
	} else {
		indexLogout(ctx, w, r)
	}
}

func indexLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, user *model.User) {
	t := template.Must(template.ParseFiles(
		TemplatePath("_base.html"),
		TemplatePath("indexLogin.html")),
	)
	type Data struct {
		commonTemplateData
		Teachers []*model.Teacher
	}
	data := &Data{commonTemplateData: getCommonTemplateData()}

	teachers, err := model.FollowingTeacherRepo.FindTeachersByUserId(user.Id)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.Teachers = teachers

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func indexLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		TemplatePath("_base.html"),
		TemplatePath("indexLogout.html")),
	)
	type Data struct {
		commonTemplateData
	}
	data := &Data{commonTemplateData: getCommonTemplateData()}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := model.GetLoggedInUser(ctx)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie(ApiTokenCookieName)
	if err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to get token cookie"))
		return
	}
	token := cookie.Value
	cookieToDelete := &http.Cookie{
		Name:     ApiTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	}
	http.SetCookie(w, cookieToDelete)
	if err := model.UserApiTokenRepo.DeleteByUserIdAndToken(user.Id, token); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
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
	fmt.Printf("!!! teacherIds = %v\n", teacherIds)
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
