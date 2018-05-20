package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/event_logger"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/interfaces/http/flash_message"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
)

const (
	followedMessage                = "フォローしました！"
	unfollowedMessage              = "削除しました！"
	emptyTeacherURLMessage         = "講師のURLまたはIDを入力して下さい"
	invalidTeacherURLMessage       = "正しい講師のURLまたはIDを入力して下さい"
	reachedMaxFollowTeacherMessage = "フォロー可能な上限数(%d)を超えました"
)

var _ = fmt.Printf

func (s *server) getMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getMe(w, r)
	}
}

func (s *server) getMe(w http.ResponseWriter, r *http.Request) {
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
		commonTemplateData: s.getCommonTemplateData(r, true, user.ID),
	}
	data.ShowTutorial = !user.FollowedTeacherAt.Valid

	mPlanService := model.NewMPlanService(s.db)
	plan, err := mPlanService.FindByPK(user.PlanID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.MPlan = plan

	followingTeacherService := model.NewFollowingTeacherService(s.db)
	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	data.Teachers = teachers

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		))
		return
	}
}

func (s *server) postMeFollowingTeachersCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.postMeFollowingTeachersCreate(w, r)
	}
}

func (s *server) postMeFollowingTeachersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	teacherIDsOrUrl := r.FormValue("teacherIdsOrUrl")
	if teacherIDsOrUrl == "" {
		e := flash_message.New(flash_message.KindWarning, emptyTeacherURLMessage)
		if err := s.flashMessageStore.Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	teachers, err := model.NewTeachersFromIDsOrURL(teacherIDsOrUrl)
	if err != nil {
		e := flash_message.New(flash_message.KindWarning, invalidTeacherURLMessage)
		if err := s.flashMessageStore.Save(e); err != nil {
			InternalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	followingTeacherService := model.NewFollowingTeacherService(s.db)
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
		if err := s.flashMessageStore.Save(e); err != nil {
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
		userService := model.NewUserService(s.db)
		if err := userService.UpdateFollowedTeacherAt(user); err != nil {
			InternalServerError(w, err)
			return
		}
		updateFollowedTeacherAt = true
	}

	mCountryService := model.NewMCountryService(s.db)
	// TODO: preload
	mCountries, err := mCountryService.LoadAll()
	if err != nil {
		InternalServerError(w, err)
		return
	}
	fetcher := fetcher.NewLessonFetcher(nil, 1, false, mCountries, logger.App)
	defer fetcher.Close()
	now := time.Now().UTC()
	teacherIDs := make([]string, 0, len(teachers))
	for _, t := range teachers {
		teacher, _, err := fetcher.Fetch(t.ID)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			InternalServerError(w, err)
			return
		}

		if _, err := followingTeacherService.FollowTeacher(user.ID, teacher, now); err != nil {
			InternalServerError(w, err)
			return
		}
		teacherIDs = append(teacherIDs, fmt.Sprint(t.ID))
	}

	go event_logger.SendGAMeasurementEvent2(
		event_logger.MustGAMeasurementEventValues(r.Context()),
		event_logger.CategoryFollowingTeacher, "follow",
		strings.Join(teacherIDs, ","), int64(len(teacherIDs)), user.ID,
	)
	if updateFollowedTeacherAt {
		go event_logger.SendGAMeasurementEvent2(
			event_logger.MustGAMeasurementEventValues(r.Context()),
			event_logger.CategoryUser, "followFirstTime",
			fmt.Sprint(user.ID), 0, user.ID,
		)
	}

	successMessage := flash_message.New(flash_message.KindSuccess, followedMessage)
	if err := s.flashMessageStore.Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
}

func (s *server) postMeFollowingTeachersDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.postMeFollowingTeachersDelete(w, r)
	}
}

func (s *server) postMeFollowingTeachersDelete(w http.ResponseWriter, r *http.Request) {
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

	followingTeacherService := model.NewFollowingTeacherService(s.db)
	_, err := followingTeacherService.DeleteTeachersByUserIDAndTeacherIDs(
		user.ID,
		util.StringToUint32Slice(teacherIDs...),
	)
	if err != nil {
		InternalServerError(w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to delete teachers"),
			errors.WithResource(errors.NewResource("following_teacher_service", "teacherIDs", teacherIDs)),
		))
		return
	}
	go event_logger.SendGAMeasurementEvent2(
		event_logger.MustGAMeasurementEventValues(r.Context()),
		event_logger.CategoryFollowingTeacher, "unfollow",
		strings.Join(teacherIDs, ","), int64(len(teacherIDs)), user.ID,
	)

	successMessage := flash_message.New(flash_message.KindSuccess, unfollowedMessage)
	if err := s.flashMessageStore.Save(successMessage); err != nil {
		InternalServerError(w, err)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
}

func (s *server) getMeSettingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getMeSetting(w, r)
	}
}

func (s *server) getMeSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/setting.html"))
	type Data struct {
		commonTemplateData
		Email string
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, true, user.ID),
		Email:              user.Email,
	}

	if err := t.Execute(w, data); err != nil {
		InternalServerError(w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		))
		return
	}
}

func (s *server) getMeLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getMeLogout(w, r)
	}
}

func (s *server) getMeLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cookie, err := r.Cookie(APITokenCookieName)
	if err != nil {
		InternalServerError(w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to get token cookie"),
		))
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
	userAPITokenService := model.NewUserAPITokenService(s.db)
	if err := userAPITokenService.DeleteByUserIDAndToken(user.ID, token); err != nil {
		InternalServerError(w, err)
		return
	}
	go event_logger.SendGAMeasurementEvent2(
		event_logger.MustGAMeasurementEventValues(r.Context()),
		event_logger.CategoryUser, "logout", "", 0, user.ID,
	)

	http.Redirect(w, r, "/", http.StatusFound)
}
