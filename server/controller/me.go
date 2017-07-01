package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/context_data"
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
	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/index.html"))
	type Data struct {
		commonTemplateData
		ShowTutorial bool
		Teachers     []*model.Teacher
		MPlan        *model.MPlan
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, true, user.ID),
	}
	data.ShowTutorial = !user.FollowedTeacherAt.Valid

	db := context_data.MustDB(ctx)
	mPlanService := model.NewMPlanService(db)
	plan, err := mPlanService.FindByPK(user.PlanID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.MPlan = plan

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
	user := context_data.MustLoggedInUser(ctx)
	teacherIDsOrUrl := r.FormValue("teacherIdsOrUrl")
	if teacherIDsOrUrl == "" {
		e := flash_message.New(flash_message.KindWarning, emptyTeacherURLMessage)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	teachers, err := model.NewTeachersFromIDsOrURL(teacherIDsOrUrl)
	if err != nil {
		e := flash_message.New(flash_message.KindWarning, invalidTeacherURLMessage)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	db := context_data.MustDB(ctx)
	followingTeacherService := model.NewFollowingTeacherService(db)
	reachesLimit, err := followingTeacherService.ReachesFollowingTeacherLimit(user.ID, len(teachers))
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if reachesLimit {
		e := flash_message.New(
			flash_message.KindWarning,
			fmt.Sprintf(reachedMaxFollowTeacherMessage, model.MaxFollowTeacherCount),
		)
		if err := flash_message.MustStore(ctx).Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	// Update user.followed_teacher_at when first time to follow teachers.
	// the column is used for showing tutorial or not.
	updateFollowedTeacherAt := false
	if !user.FollowedTeacherAt.Valid {
		userService := model.NewUserService(db)
		if err := userService.UpdateFollowedTeacherAt(user); err != nil {
			InternalServerError(w, err)
			return
		}
		updateFollowedTeacherAt = true
	}

	mCountryService := model.NewMCountryService(db)
	// TODO: preload
	mCountries, err := mCountryService.LoadAll()
	if err != nil {
		InternalServerError(w, err)
		return
	}
	fetcher := fetcher.NewTeacherLessonFetcher(nil, 1, false, mCountries, logger.App)
	now := time.Now().UTC()
	teacherIDs := make([]string, 0, len(teachers))
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
		if _, err := followingTeacherService.FollowTeacher(user.ID, teacher, now); err != nil {
			InternalServerError(w, err)
			return
		}
		teacherIDs = append(teacherIDs, fmt.Sprint(t.ID))
	}

	go sendMeasurementEvent(
		r, eventCategoryFollowingTeacher, "follow",
		strings.Join(teacherIDs, ","), int64(len(teacherIDs)), user.ID,
	)
	if updateFollowedTeacherAt {
		go sendMeasurementEvent(
			r, eventCategoryUser, "followFirstTime",
			fmt.Sprint(user.ID), 0, user.ID,
		)
	}

	successMessage := flash_message.New(flash_message.KindSuccess, followedMessage)
	if err := flash_message.MustStore(ctx).Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
}

func PostMeFollowingTeachersDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	teacherIDs := r.Form["teacherIds"]
	if len(teacherIDs) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	followingTeacherService := model.NewFollowingTeacherService(context_data.MustDB(ctx))
	_, err := followingTeacherService.DeleteTeachersByUserIDAndTeacherIDs(
		user.ID,
		util.StringToUint32Slice(teacherIDs...),
	)
	if err != nil {
		e := errors.InternalWrapf(err, "Failed to delete Teachers: teacherIds=%v", teacherIDs)
		InternalServerError(w, e)
		return
	}
	go sendMeasurementEvent(
		r, eventCategoryFollowingTeacher, "unfollow",
		strings.Join(teacherIDs, ","), int64(len(teacherIDs)), user.ID,
	)

	successMessage := flash_message.New(flash_message.KindSuccess, unfollowedMessage)
	if err := flash_message.MustStore(ctx).Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
}

func GetMeSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/setting.html"))
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: getCommonTemplateData(r, true, user.ID),
		Email:              user.Email,
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to template.Execute()"))
		return
	}
}

func PostMeSettingUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	email := r.FormValue("email")
	// TODO: better validation
	if email == "" || !validateEmail(email) {
		http.Redirect(w, r, "/me/setting", http.StatusFound)
		return
	}

	service := model.NewUserService(context_data.MustDB(ctx))
	if err := service.UpdateEmail(user, email); err != nil {
		InternalServerError(w, err)
		return
	}
	go sendMeasurementEvent(r, eventCategoryUser, "update", fmt.Sprint(user.ID), 0, user.ID)

	successMessage := flash_message.New(flash_message.KindSuccess, updatedMessage)
	if err := flash_message.MustStore(ctx).Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me/setting?"+successMessage.AsURLQueryString(), http.StatusFound)
}

func GetMeLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie(APITokenCookieName)
	if err != nil {
		InternalServerError(w, errors.InternalWrapf(err, "Failed to get token cookie"))
		return
	}
	token := cookie.Value
	cookieToDelete := &http.Cookie{
		Name:     APITokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	}
	http.SetCookie(w, cookieToDelete)
	userAPITokenService := model.NewUserAPITokenService(context_data.MustDB(ctx))
	if err := userAPITokenService.DeleteByUserIDAndToken(user.ID, token); err != nil {
		InternalServerError(w, err)
		return
	}
	go sendMeasurementEvent(r, eventCategoryUser, "logout", "", 0, user.ID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func validateEmail(email string) bool {
	// TODO: better validation
	return strings.Contains(email, "@")
}
