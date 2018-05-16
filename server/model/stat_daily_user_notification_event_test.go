package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatDailyUserNotificationEventService_CreateOrUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	user := helper.CreateUser("test", "test@gmail.com")
	datetime := time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC)
	const num = 3
	err := createEventLogEmails(user.ID, datetime, num)
	r.NoError(err)

	err = statDailyUserNotificationEventService.CreateOrUpdate(datetime)
	r.NoError(err)
	events, err := statDailyUserNotificationEventService.FindAllByDate(datetime)
	r.NoError(err)
	a.Equal(1, len(events))
	a.Equal(user.ID, events[0].UserID)
	a.EqualValues(num, events[0].Count)
}

func createEventLogEmails(userID uint32, datetime time.Time, num int) error {
	for i := 0; i < num; i++ {
		err := eventLogEmailService.Create(&EventLogEmail{
			Datetime:   datetime,
			Event:      "open",
			EmailType:  EmailTypeNewLessonNotifier,
			UserID:     userID,
			UserAgent:  "",
			TeacherIDs: "",
			URL:        "",
		})
		if err != nil {
			return err
		}
	}
	return nil
}
