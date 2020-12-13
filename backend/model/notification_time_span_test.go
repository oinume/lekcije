package model

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestNotificationTimeSpan_Within(t *testing.T) {
	table := []struct {
		name     string
		fromTime string
		toTime   string
		target   string
		within   bool
	}{
		{
			name:     "normal true",
			fromTime: "08:30:00",
			toTime:   "09:30:00",
			target:   "09:30:00",
			within:   true,
		},
		{
			name:     "normal false",
			fromTime: "08:30:00",
			toTime:   "09:30:00",
			target:   "10:00:00",
			within:   false,
		},
		{
			name:     "fromTime > toTime",
			fromTime: "23:00:00",
			toTime:   "01:00:00",
			target:   "23:30:00",
			within:   true,
		},
	}

	for _, tt := range table {
		timeSpan := NotificationTimeSpan{
			Number:   1,
			FromTime: tt.fromTime,
			ToTime:   tt.toTime,
		}
		if err := timeSpan.ParseTime(); err != nil {
			t.Fatalf("timeSpan.ParseTime() failed: err=%v", err)
		}
		target, _ := time.Parse("15:04:05", tt.target)
		if got := timeSpan.Within(target); got != tt.within {
			t.Errorf("%v: unpexpected value from Within: target=%v, fromTime=%v, toTime=%v: got=%v, want=%v",
				tt.name, tt.target, tt.fromTime, tt.toTime, got, tt.within,
			)
		}
	}
}

func TestNotificationTimeSpanList_Within(t *testing.T) {
	table := []struct {
		name         string
		timeSpanList NotificationTimeSpanList
		target       string
		within       bool
	}{
		{
			name: "normal true",
			timeSpanList: NotificationTimeSpanList{
				&NotificationTimeSpan{
					Number:   1,
					FromTime: "01:00:00",
					ToTime:   "02:00:00",
				},
				&NotificationTimeSpan{
					Number:   2,
					FromTime: "08:30:00",
					ToTime:   "10:00:00",
				},
			},
			target: "09:00:00",
			within: true,
		},
		{
			name: "normal false",
			timeSpanList: NotificationTimeSpanList{
				&NotificationTimeSpan{
					Number:   1,
					FromTime: "01:00:00",
					ToTime:   "02:00:00",
				},
				&NotificationTimeSpan{
					Number:   2,
					FromTime: "08:30:00",
					ToTime:   "10:00:00",
				},
			},
			target: "10:30:00",
			within: false,
		},
	}

	for _, tt := range table {
		target, _ := time.Parse("15:04:05", tt.target)
		if got := tt.timeSpanList.Within(target); got != tt.within {
			t.Errorf("%v: unexpected value from Within: got=%v, want=%v", tt.name, got, tt.within)
		}
	}
}

func TestNotificationTimeSpanService_UpdateAll(t *testing.T) {
	user := helper.CreateRandomUser(t)
	now := time.Now().UTC()
	timeSpans := []*NotificationTimeSpan{
		{
			UserID:    user.ID,
			Number:    1,
			FromTime:  "08:30:00",
			ToTime:    "09:30:00",
			CreatedAt: now,
		},
		{
			UserID:    user.ID,
			Number:    2,
			FromTime:  "20:30:00",
			ToTime:    "21:30:00",
			CreatedAt: now,
		},
	}

	err := notificationTimeSpanService.UpdateAll(user.ID, timeSpans)
	if err != nil {
		t.Fatalf("notificationTimeSpanService.UpdateAll failed: err=%v", err)
	}

	gotTimeSpans, err := notificationTimeSpanService.FindByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("notificationTimeSpanService.FindByUserID failed: err=%v", err)
	}
	if got, want := len(gotTimeSpans), len(timeSpans); got != want {
		t.Fatalf("length doesn't match : got=%v, want=%v", got, want)
	}
	for i, got := range gotTimeSpans {
		// Maybe gorm overwrite CreatedAt somehow
		timeSpans[i].CreatedAt = time.Time{}
		got.CreatedAt = time.Time{}
		if want := timeSpans[i]; !reflect.DeepEqual(got, want) {
			t.Errorf("unexepected value[%v]: \ngot =%v\nwant=%v", i, got, want)
		}
	}
}
