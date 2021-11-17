package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/fetcher"
	"github.com/oinume/lekcije/backend/interface/http/flash_message"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/util"
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
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	data.MPlan = plan

	followingTeacherService := model.NewFollowingTeacherService(s.db)
	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
	if err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	data.Teachers = teachers

	if err := t.Execute(w, data); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), user.ID)
		return
	}
}

func (s *server) getMeNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	t := ParseHTMLTemplates(TemplatePath("me/new.html"))
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
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	data.MPlan = plan

	followingTeacherService := model.NewFollowingTeacherService(s.db)
	teachers, err := followingTeacherService.FindTeachersByUserID(user.ID)
	if err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	data.Teachers = teachers

	if err := t.Execute(w, data); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), user.ID)
		return
	}
}

func (s *server) postMeFollowingTeachersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context_data.MustLoggedInUser(ctx)
	teacherIDsOrURL := r.FormValue("teacherIdsOrUrl")
	if teacherIDsOrURL == "" {
		e := flash_message.New(flash_message.KindWarning, emptyTeacherURLMessage)
		if err := s.flashMessageStore.Save(e); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	teachers, err := model.NewTeachersFromIDsOrURL(teacherIDsOrURL)
	if err != nil {
		e := flash_message.New(flash_message.KindWarning, invalidTeacherURLMessage)
		if err := s.flashMessageStore.Save(e); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}
		http.Redirect(w, r, "/me?"+e.AsURLQueryString(), http.StatusFound)
		return
	}

	followingTeacherService := model.NewFollowingTeacherService(s.db)
	reachesLimit, err := followingTeacherService.ReachesFollowingTeacherLimit(user.ID, len(teachers))
	if err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	if reachesLimit {
		e := flash_message.New(
			flash_message.KindWarning,
			fmt.Sprintf(reachedMaxFollowTeacherMessage, model.MaxFollowTeacherCount),
		)
		if err := s.flashMessageStore.Save(e); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
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
		// TODO: 1回でまとめて更新する
		if err := userService.UpdateFollowedTeacherAt(user); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}
		if err := userService.UpdateOpenNotificationAt(user.ID, time.Now().UTC()); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}
		updateFollowedTeacherAt = true
	}

	mCountryService := model.NewMCountryService(s.db)
	// TODO: preload
	mCountries, err := mCountryService.LoadAll(ctx)
	if err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	fetcher := fetcher.NewLessonFetcher(nil, 1, false, mCountries, s.appLogger)
	defer fetcher.Close()
	now := time.Now().UTC()
	teacherIDs := make([]string, 0, len(teachers))
	for _, t := range teachers {
		teacher, _, err := fetcher.Fetch(r.Context(), t.ID)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}

		if _, err := followingTeacherService.FollowTeacher(user.ID, teacher, now); err != nil {
			internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
			return
		}
		teacherIDs = append(teacherIDs, fmt.Sprint(t.ID))
	}

	go s.sendGAMeasurementEvent(
		r.Context(),
		model2.GAMeasurementEventCategoryFollowingTeacher,
		"follow",
		strings.Join(teacherIDs, ","),
		int64(len(teacherIDs)),
		user.ID,
	)

	if updateFollowedTeacherAt {
		go s.sendGAMeasurementEvent(
			r.Context(),
			model2.GAMeasurementEventCategoryUser,
			"followFirstTime",
			fmt.Sprint(user.ID),
			0,
			user.ID,
		)
	}

	successMessage := flash_message.New(flash_message.KindSuccess, followedMessage)
	if err := s.flashMessageStore.Save(successMessage); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
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
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to delete teachers"),
			errors.WithResource(errors.NewResource("following_teacher_service", "teacherIDs", teacherIDs)),
		), user.ID)
		return
	}

	go s.sendGAMeasurementEvent(
		r.Context(),
		model2.GAMeasurementEventCategoryUser,
		"unfollow",
		strings.Join(teacherIDs, ","),
		int64(len(teacherIDs)),
		user.ID,
	)

	successMessage := flash_message.New(flash_message.KindSuccess, unfollowedMessage)
	if err := s.flashMessageStore.Save(successMessage); err != nil {
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}

	http.Redirect(w, r, "/me?"+successMessage.AsURLQueryString(), http.StatusFound)
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
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), user.ID)
		return
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
		internalServerError(r.Context(), s.errorRecorder, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to get token cookie"),
		), user.ID)
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
		internalServerError(r.Context(), s.errorRecorder, w, err, user.ID)
		return
	}
	go s.sendGAMeasurementEvent(
		r.Context(),
		model2.GAMeasurementEventCategoryUser,
		"logout",
		fmt.Sprint(user.ID),
		0,
		user.ID,
	)
	http.Redirect(w, r, "/", http.StatusFound)
}
