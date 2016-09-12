package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"golang.org/x/net/context"
)

var _ = fmt.Printf

func PostMeFollowingTeachersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	teacherIDsOrUrl := r.FormValue("teacherIdsOrUrl")
	if teacherIDsOrUrl == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	teachers, err := model.NewTeachersFromIDsOrURL(teacherIDsOrUrl)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	fetcher := fetcher.NewTeacherLessonFetcher(http.DefaultClient, logger.AppLogger)

	for _, t := range teachers {
		teacher, _, err := fetcher.Fetch(t.ID)
		if err != nil {
			// TODO: continue the loop
			InternalServerError(w, err)
			return
		}

		now := time.Now()
		teacher.CreatedAt = now
		teacher.UpdatedAt = now
		// TODO: Create method on service (FollowTeacher)
		db := model.MustDB(ctx)
		if err := db.FirstOrCreate(teacher).Error; err != nil {
			e := errors.InternalWrapf(err, "Failed to create Teacher: teacherID=%d", teacher.ID)
			InternalServerError(w, e)
			return
		}

		ft := &model.FollowingTeacher{
			UserID:    user.ID,
			TeacherID: teacher.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := db.FirstOrCreate(ft).Error; err != nil {
			e := errors.InternalWrapf(
				err,
				"Failed to create FollowingTeacher: userID=%d, teacherID=%d",
				user.ID, teacher.ID,
			)
			InternalServerError(w, e)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func PostMeFollowingTeachersDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	teacherIDs := r.Form["teacherIds"]
	if len(teacherIDs) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	followingTeacherService := model.NewFollowingTeacherService(model.MustDB(ctx))
	_, err := followingTeacherService.DeleteTeachersByUserIDAndTeacherIDs(
		user.ID,
		util.StringToUint32Slice(teacherIDs...),
	)
	if err != nil {
		e := errors.InternalWrapf(err, "Failed to delete Teachers: teacherIds=%v", teacherIDs)
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
		Email:              user.Email.Raw(),
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func PostMeSettingUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := model.MustLoggedInUser(ctx)
	email := r.FormValue("email")
	// TODO: better validation
	if email == "" || !validateEmail(email) {
		http.Redirect(w, r, "/me/setting", http.StatusFound)
		return
	}

	service := model.NewUserService(model.MustDB(ctx))
	if err := service.UpdateEmail(user, email); err != nil {
		InternalServerError(w, err)
		return
	}
	http.Redirect(w, r, "/me/setting", http.StatusFound)
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}
