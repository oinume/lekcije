package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/backend/model2"
)

func TestStatDailyUserNotificationEventService_CreateOrUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	helper.TruncateAllTables(t)

	user1 := helper.CreateUser(t, "test1", "test1@gmail.com")
	user2 := helper.CreateUser(t, "test2", "test2@gmail.com")
	testCases := []struct {
		userID   uint32
		datetime time.Time
		count    int
	}{
		{
			userID:   user1.ID,
			datetime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			count:    3,
		},
		{
			userID:   user2.ID,
			datetime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			count:    0,
		},
	}

	for i, tc := range testCases {
		err := createEventLogEmails(tc.userID, tc.datetime, tc.count)
		r.NoError(err)
		err = statDailyUserNotificationEventService.CreateOrUpdate(tc.datetime)
		r.NoError(err)
		events, err := statDailyUserNotificationEventService.FindAllByDate(tc.datetime)
		r.NoError(err)

		r.Equal(len(testCases), len(events))
		a.Equal(tc.userID, events[i].UserID)
		a.Equal("open", events[i].Event)
		a.EqualValues(tc.count, events[i].Count)
	}
}

func createEventLogEmails(userID uint32, datetime time.Time, num int) error {
	for i := 0; i < num; i++ {
		err := eventLogEmailService.Create(&EventLogEmail{
			Datetime:   datetime.Add(time.Duration(i) * time.Second),
			Event:      "open",
			EmailType:  model2.EmailTypeNewLessonNotifier,
			UserID:     userID,
			UserAgent:  "test",
			TeacherIDs: "1",
			URL:        "",
		})
		if err != nil {
			return err
		}
	}
	return nil
}
