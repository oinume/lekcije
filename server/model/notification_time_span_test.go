package model

import (
	"reflect"
	"testing"
	"time"
)

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

	err := notificationTimeSpanService.UpdateAll(timeSpans)
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
