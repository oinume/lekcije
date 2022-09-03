// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package graphqltest

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

type NotificationTimeSpanInput struct {
	FromHour   int `json:"fromHour"`
	FromMinute int `json:"fromMinute"`
	ToHour     int `json:"toHour"`
	ToMinute   int `json:"toMinute"`
}

// GetFromHour returns NotificationTimeSpanInput.FromHour, and is useful for accessing the field via an interface.
func (v *NotificationTimeSpanInput) GetFromHour() int { return v.FromHour }

// GetFromMinute returns NotificationTimeSpanInput.FromMinute, and is useful for accessing the field via an interface.
func (v *NotificationTimeSpanInput) GetFromMinute() int { return v.FromMinute }

// GetToHour returns NotificationTimeSpanInput.ToHour, and is useful for accessing the field via an interface.
func (v *NotificationTimeSpanInput) GetToHour() int { return v.ToHour }

// GetToMinute returns NotificationTimeSpanInput.ToMinute, and is useful for accessing the field via an interface.
func (v *NotificationTimeSpanInput) GetToMinute() int { return v.ToMinute }

type UpdateNotificationTimeSpansInput struct {
	TimeSpans []NotificationTimeSpanInput `json:"timeSpans"`
}

// GetTimeSpans returns UpdateNotificationTimeSpansInput.TimeSpans, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansInput) GetTimeSpans() []NotificationTimeSpanInput {
	return v.TimeSpans
}

// UpdateNotificationTimeSpansResponse is returned by UpdateNotificationTimeSpans on success.
type UpdateNotificationTimeSpansResponse struct {
	UpdateNotificationTimeSpans UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload `json:"updateNotificationTimeSpans"`
}

// GetUpdateNotificationTimeSpans returns UpdateNotificationTimeSpansResponse.UpdateNotificationTimeSpans, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansResponse) GetUpdateNotificationTimeSpans() UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload {
	return v.UpdateNotificationTimeSpans
}

// UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload includes the requested fields of the GraphQL type NotificationTimeSpanPayload.
type UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload struct {
	TimeSpans []UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan `json:"timeSpans"`
}

// GetTimeSpans returns UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload.TimeSpans, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload) GetTimeSpans() []UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan {
	return v.TimeSpans
}

// UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan includes the requested fields of the GraphQL type NotificationTimeSpan.
type UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan struct {
	FromHour   int `json:"fromHour"`
	FromMinute int `json:"fromMinute"`
	ToHour     int `json:"toHour"`
	ToMinute   int `json:"toMinute"`
}

// GetFromHour returns UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan.FromHour, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan) GetFromHour() int {
	return v.FromHour
}

// GetFromMinute returns UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan.FromMinute, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan) GetFromMinute() int {
	return v.FromMinute
}

// GetToHour returns UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan.ToHour, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan) GetToHour() int {
	return v.ToHour
}

// GetToMinute returns UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan.ToMinute, and is useful for accessing the field via an interface.
func (v *UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan) GetToMinute() int {
	return v.ToMinute
}

type UpdateViewerInput struct {
	Email string `json:"email"`
}

// GetEmail returns UpdateViewerInput.Email, and is useful for accessing the field via an interface.
func (v *UpdateViewerInput) GetEmail() string { return v.Email }

// UpdateViewerResponse is returned by UpdateViewer on success.
type UpdateViewerResponse struct {
	UpdateViewer UpdateViewerUpdateViewerUser `json:"updateViewer"`
}

// GetUpdateViewer returns UpdateViewerResponse.UpdateViewer, and is useful for accessing the field via an interface.
func (v *UpdateViewerResponse) GetUpdateViewer() UpdateViewerUpdateViewerUser { return v.UpdateViewer }

// UpdateViewerUpdateViewerUser includes the requested fields of the GraphQL type User.
type UpdateViewerUpdateViewerUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// GetId returns UpdateViewerUpdateViewerUser.Id, and is useful for accessing the field via an interface.
func (v *UpdateViewerUpdateViewerUser) GetId() string { return v.Id }

// GetEmail returns UpdateViewerUpdateViewerUser.Email, and is useful for accessing the field via an interface.
func (v *UpdateViewerUpdateViewerUser) GetEmail() string { return v.Email }

// __UpdateNotificationTimeSpansInput is used internally by genqlient
type __UpdateNotificationTimeSpansInput struct {
	Input UpdateNotificationTimeSpansInput `json:"input"`
}

// GetInput returns __UpdateNotificationTimeSpansInput.Input, and is useful for accessing the field via an interface.
func (v *__UpdateNotificationTimeSpansInput) GetInput() UpdateNotificationTimeSpansInput {
	return v.Input
}

// __UpdateViewerInput is used internally by genqlient
type __UpdateViewerInput struct {
	Input UpdateViewerInput `json:"input"`
}

// GetInput returns __UpdateViewerInput.Input, and is useful for accessing the field via an interface.
func (v *__UpdateViewerInput) GetInput() UpdateViewerInput { return v.Input }

func UpdateNotificationTimeSpans(
	ctx context.Context,
	client graphql.Client,
	input UpdateNotificationTimeSpansInput,
) (*UpdateNotificationTimeSpansResponse, error) {
	req := &graphql.Request{
		OpName: "UpdateNotificationTimeSpans",
		Query: `
mutation UpdateNotificationTimeSpans ($input: UpdateNotificationTimeSpansInput!) {
	updateNotificationTimeSpans(input: $input) {
		timeSpans {
			fromHour
			fromMinute
			toHour
			toMinute
		}
	}
}
`,
		Variables: &__UpdateNotificationTimeSpansInput{
			Input: input,
		},
	}
	var err error

	var data UpdateNotificationTimeSpansResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func UpdateViewer(
	ctx context.Context,
	client graphql.Client,
	input UpdateViewerInput,
) (*UpdateViewerResponse, error) {
	req := &graphql.Request{
		OpName: "UpdateViewer",
		Query: `
mutation UpdateViewer ($input: UpdateViewerInput!) {
	updateViewer(input: $input) {
		id
		email
	}
}
`,
		Variables: &__UpdateViewerInput{
			Input: input,
		},
	}
	var err error

	var data UpdateViewerResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
