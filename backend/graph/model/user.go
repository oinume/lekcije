package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	// NOTE: FollowingTeachers is not defined intentionally because it is provided by userResolver
}
