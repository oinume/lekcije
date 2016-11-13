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
	followedMessage                = "フォローしました！"
	unfollowedMessage              = "削除しました！"
	updatedMessage                 = "設定を更新しました！"
	emptyTeacherURLMessage         = "講師のURLまたはIDを入力して下さい"
	invalidTeacherURLMessage       = "正しい講師のURLまたはIDを入力して下さい"
	reachedMaxFollowTeacherMessage = "フォロー可能な上限数(%d)を超えました"
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
	data.ShowTutorial = !user.FollowedTeacherAt.Valid

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
		e := flash_message.New(flash_message.KindWarning, emptyTeacherURLMessage)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?flashMessageKey="+e.Key, http.StatusFound)
		return
	}

	teachers, err := model.NewTeachersFromIDsOrURL(teacherIDsOrUrl)
	if err != nil {
		e := flash_message.New(flash_message.KindWarning, invalidTeacherURLMessage)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?flashMessageKey="+e.Key, http.StatusFound)
		return
	}

	db := model.MustDB(ctx)
	followingTeacherService := model.NewFollowingTeacherService(db)
	count, err := followingTeacherService.CountFollowingTeachersByUserID(user.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if count+len(teachers) > model.MaxFollowTeacherCount { // TODO: As plan
		e := flash_message.New(
			flash_message.KindWarning,
			fmt.Sprintf(reachedMaxFollowTeacherMessage, model.MaxFollowTeacherCount),
		)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		// TODO: Generate flashMessageKey=fugafuga by e.ToURLParam
		http.Redirect(w, r, "/me?flashMessageKey="+e.Key, http.StatusFound)
		return
	}

	// Update user.followed_teacher_at when first time to follow teachers.
	// user.followed_teacher_at is used for showing tutorial.
	if !user.FollowedTeacherAt.Valid {
		userService := model.NewUserService(db)
		if err := userService.UpdateFollowedTeacherAt(user); err != nil {
			InternalServerError(w, err)
			return
		}
	}

	fetcher := fetcher.NewTeacherLessonFetcher(nil, logger.AppLogger)
	now := time.Now().UTC()
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
		if err := followingTeacherService.FollowTeacher(user.ID, teacher, now); err != nil {
			InternalServerError(w, err)
			return
		}
	}

	successMessage := flash_message.New(flash_message.KindSuccess, followedMessage)
	if err := flash_message.MustStore(ctx).Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me?flashMessageKey="+successMessage.Key, http.StatusFound)
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
