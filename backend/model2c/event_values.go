package model2

type GAMeasurementEvent struct {
	UserAgentOverride string
	ClientID          string
	DocumentHostName  string
	DocumentPath      string
	DocumentTitle     string
	DocumentReferrer  string
	IPOverride        string
}

const (
	GAMeasurementEventCategoryEmail            = "email"
	GAMeasurementEventCategoryUser             = "user"
	GAMeasurementEventCategoryFollowingTeacher = "followingTeacher"
)
