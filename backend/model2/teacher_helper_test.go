package model2_test

import (
	"fmt"
	"testing"

	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/model2"
)

func TestNewTeacherFromIDOrURL(t *testing.T) {
	tests := map[string]struct {
		idOrURL string
		want    *model2.Teacher
		wantErr error
	}{
		"id": {
			idOrURL: "12345",
			want: &model2.Teacher{
				ID: 12345,
			},
		},
		"url": {
			idOrURL: "https://eikaiwa.dmm.com/teacher/index/67890/",
			want: &model2.Teacher{
				ID: 67890,
			},
		},
		"invalid url": {
			idOrURL: "https://www.dmm.com/teacher/index/67890/",
			wantErr: fmt.Errorf("code.InvalidArgument: Failed to parse idOrURL: "),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := model2.NewTeacherFromIDOrURL(test.idOrURL)
			if err != nil {
				assertion.RequireEqual(t, test.wantErr.Error(), err.Error(), "unexpected error")
			}
			assertion.AssertEqual(t, test.want, got, "unexpected teacher")
		})
	}
}
