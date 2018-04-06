package model

import (
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
		want     bool
	}{
		{
			name:     "normal true",
			fromTime: "08:30:00",
			toTime:   "09:30:00",
			target:   "09:30:00",
			want:     true,
		},
		{
			name:     "normal false",
			fromTime: "08:30:00",
			toTime:   "09:30:00",
			target:   "10:00:00",
			want:     false,
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
		if got := timeSpan.Within(target); got != tt.want {
			t.Errorf("%v: unpexpected value from Within: got=%v, want=%v", tt.name, got, tt.want)
		}
	}
}

func TestNotificationTimeSpanList_Within(t *testing.T) {
	list := []*NotificationTimeSpanList{
		&NotificationTimeSpan{
			Number: 1,
		},
	}
}

func TestNotificationTimeSpanService_UpdateAll(t *testing.T) {
	user := helper.CreateRandomUser()
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

	gotTimeSpans, err := notificationTimeSpanService.FindByUserID(user.ID)
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
