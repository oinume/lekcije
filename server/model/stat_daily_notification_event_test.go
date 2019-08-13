package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatDailyNotificationEventService_CreateOrUpdate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	helper.TruncateAllTables(t)

	user1 := helper.CreateUser(t, "test1", "test1@gmail.com")
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
	}

	for i, tc := range testCases {
		err := createEventLogEmails(tc.userID, tc.datetime, tc.count)
		r.NoError(err)
		stats, err := eventLogEmailService.FindStatDailyNotificationEventByDate(tc.datetime)
		r.NoError(err)
		statUUs, err := eventLogEmailService.FindStatDailyNotificationEventUUCountByDate(tc.datetime)
		r.NoError(err)

		values := make(map[string]*StatDailyNotificationEvent, 100)
		for _, s := range stats {
			values[s.Event] = s
		}
		for _, s := range statUUs {
			v := values[s.Event]
			v.UUCount = s.UUCount
			if err := statDailyNotificationEventService.CreateOrUpdate(v); err != nil {
				r.NoError(err)
			}
		}

		events, err := statDailyNotificationEventService.FindAllByDate(tc.datetime)
		r.NoError(err)

		r.Equal(len(testCases), len(events))
		a.Equal("open", events[i].Event)
		a.EqualValues(tc.count, events[i].Count)
	}
}
