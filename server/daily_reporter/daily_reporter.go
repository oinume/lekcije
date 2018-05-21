package daily_reporter

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/model"
)

type Main struct {
	TargetDate *string
	LogLevel   *string
	DB         *gorm.DB
}

func (m *Main) Run() error {
	if *m.TargetDate == "" {
		return fmt.Errorf("-target-date is required")
	}
	date, err := time.Parse("2006-01-02", *m.TargetDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", *m.TargetDate)
	}

	if err := m.createStatNewLessonNotifier(date); err != nil {
		return err
	}
	if err := m.createStatDailyUserNotificationEvent(date); err != nil {
		return err
	}
	return nil
}

func (m *Main) createStatNewLessonNotifier(date time.Time) error {
	service := model.NewEventLogEmailService(m.DB)
	stats, err := service.FindStatDailyNotificationEventByDate(date)
	if err != nil {
		return err
	}
	statUUs, err := service.FindStatDailyNotificationEventUUCountByDate(date)
	if err != nil {
		return err
	}

	values := make(map[string]*model.StatDailyNotificationEvent, 100)
	for _, s := range stats {
		values[s.Event] = s
	}

	statDailyNotificationEventService := model.NewStatDailyNotificationEventService(m.DB)
	for _, s := range statUUs {
		v := values[s.Event]
		v.UUCount = s.UUCount
		if err := statDailyNotificationEventService.CreateOrUpdate(v); err != nil {
			return err
		}
	}

	//statsNewLessonNotifierService := model.NewStatsNewLessonNotifierService(m.DB)
	//for _, s := range statUUs {
	//	v := values[s.Event]
	//	v.UUCount = s.UUCount
	//	if err := statsNewLessonNotifierService.CreateOrUpdate(v); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (m *Main) createStatDailyUserNotificationEvent(date time.Time) error {
	service := model.NewStatDailyUserNotificationEventService(m.DB)
	return service.CreateOrUpdate(date)
}
