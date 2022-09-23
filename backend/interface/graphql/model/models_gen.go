// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Connection interface {
	IsConnection()
	GetPageInfo() *PageInfo
	GetEdges() []Edge
	GetNodes() []Node
}

type Edge interface {
	IsEdge()
	GetCursor() string
	GetNode() Node
}

type Node interface {
	IsNode()
	GetID() string
}

type CreateFollowingTeacherInput struct {
	TeacherIDOrURL string `json:"teacherIdOrUrl"`
}

type CreateFollowingTeacherPayload struct {
	ID string `json:"id"`
}

type Empty struct {
	ID string `json:"id"`
}

type FollowingTeacher struct {
	ID        string   `json:"id"`
	Teacher   *Teacher `json:"teacher"`
	CreatedAt string   `json:"createdAt"`
}

func (FollowingTeacher) IsNode()            {}
func (this FollowingTeacher) GetID() string { return this.ID }

type FollowingTeacherConnection struct {
	PageInfo *PageInfo               `json:"pageInfo"`
	Edges    []*FollowingTeacherEdge `json:"edges"`
	Nodes    []*FollowingTeacher     `json:"nodes"`
}

func (FollowingTeacherConnection) IsConnection()               {}
func (this FollowingTeacherConnection) GetPageInfo() *PageInfo { return this.PageInfo }
func (this FollowingTeacherConnection) GetEdges() []Edge {
	if this.Edges == nil {
		return nil
	}
	interfaceSlice := make([]Edge, 0, len(this.Edges))
	for _, concrete := range this.Edges {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this FollowingTeacherConnection) GetNodes() []Node {
	if this.Nodes == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Nodes))
	for _, concrete := range this.Nodes {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

type FollowingTeacherEdge struct {
	Cursor string            `json:"cursor"`
	Node   *FollowingTeacher `json:"node"`
}

func (FollowingTeacherEdge) IsEdge()                {}
func (this FollowingTeacherEdge) GetCursor() string { return this.Cursor }
func (this FollowingTeacherEdge) GetNode() Node     { return *this.Node }

type NotificationTimeSpan struct {
	FromHour   int `json:"fromHour"`
	FromMinute int `json:"fromMinute"`
	ToHour     int `json:"toHour"`
	ToMinute   int `json:"toMinute"`
}

type NotificationTimeSpanInput struct {
	FromHour   int `json:"fromHour"`
	FromMinute int `json:"fromMinute"`
	ToHour     int `json:"toHour"`
	ToMinute   int `json:"toMinute"`
}

type NotificationTimeSpanPayload struct {
	TimeSpans []*NotificationTimeSpan `json:"timeSpans"`
}

type PageInfo struct {
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

type Teacher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateNotificationTimeSpansInput struct {
	TimeSpans []*NotificationTimeSpanInput `json:"timeSpans"`
}

type UpdateViewerInput struct {
	Email *string `json:"email"`
}
