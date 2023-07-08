package errors

import (
	"fmt"

	"github.com/morikuni/failure"
)

func NewUserIDContext(userID uint) failure.Context {
	return failure.Context{"userID": fmt.Sprint(userID)}
}

func WithUserIDContext(c failure.Context, userID uint) failure.Context {
	if c == nil {
		c = failure.Context{}
	}
	c["userID"] = fmt.Sprint(userID)
	return c
}

func WithTableContext(c failure.Context, table string) failure.Context {
	if c == nil {
		c = failure.Context{}
	}
	c["table"] = table
	return c
}
