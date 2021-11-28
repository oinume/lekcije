package model2

func (o *User) IsFollowedTeacher() bool {
	return o.FollowedTeacherAt.Valid && !o.FollowedTeacherAt.Time.IsZero()
}
