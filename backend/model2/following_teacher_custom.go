package model2

import "fmt"

func (o *FollowingTeacher) ID() string {
	return fmt.Sprintf("%v-%v", o.UserID, o.TeacherID)
}
