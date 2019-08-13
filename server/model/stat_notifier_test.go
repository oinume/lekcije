package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStatNotifierService_CreateOrUpdate(t *testing.T) {
	//a := assert.New(t)
	r := require.New(t)
	helper.TruncateAllTables(t)

	//user1 := helper.CreateUser("test1", "test1@gmail.com")
	statNotifier := &StatNotifier{
		Datetime:             time.Now().UTC(),
		Interval:             10,
		Elapsed:              100000,
		UserCount:            800,
		FollowedTeacherCount: 1600,
	}
	err := statNotifierService.CreateOrUpdate(statNotifier)
	r.NoError(err)
}
