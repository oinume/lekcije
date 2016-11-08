package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/controller/flash_message"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
)

const (
	followedMessage   = "フォローしました！"
	unfollowedMessage = "削除しました！"
	updatedMessage    = "設定を更新しました！"
)

var _ = fmt.Printf

func GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := model.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/index.html"))
	type Data struct {
		commonTemplateData
		ShowTutorial bool
		Teachers     []*model.Teacher
		Plan         *model.Plan
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, true),
	}

	showTutorialStr := r.FormValue("showTutorial")
	if showTutorialStr == "true" {
		data.ShowTutorial = true
	}

	db := model.MustDB(ctx)
	planService := model.NewPlanService(db)
	plan, err := planService.FindByPk(user.PlanID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.Plan = plan

	followingTeacherService := model.NewFollowingTeacherService(db)
	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
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

func PostMeFollowingTeachersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	fetcher := fetcher.NewTeacherLessonFetcher(nil, logger.AppLogger)

	for _, t := range teachers {
		teacher, _, err := fetcher.Fetch(t.ID)
		if err != nil {
			switch err.(type) {
			case *errors.NotFound:
				// TODO: return error to user
				continue
			default:
				// TODO: continue the loop
				InternalServerError(w, err)
				return
			}
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

	flashMessage := flash_message.New(flash_message.KindSuccess, followedMessage)
	flash_message.MustStore(ctx).Save(flashMessage)

	http.Redirect(w, r, "/?flashMessageKey="+flashMessage.Key, http.StatusFound)
}

func PostMeFollowingTeachersDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	flashMessage := flash_message.New(flash_message.KindSuccess, unfollowedMessage)
	flash_message.MustStore(ctx).Save(flashMessage)

	http.Redirect(w, r, "/?flashMessageKey="+flashMessage.Key, http.StatusFound)
}

func GetMeSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := model.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/setting.html"))
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, true),
		Email:              user.Email.Raw(),
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func PostMeSettingUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	flashMessage := flash_message.New(flash_message.KindSuccess, updatedMessage)
	flash_message.MustStore(ctx).Save(flashMessage)

	http.Redirect(w, r, "/me/setting?flashMessageKey="+flashMessage.Key, http.StatusFound)
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}
