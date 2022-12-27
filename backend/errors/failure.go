package errors

import (
	"fmt"

	"github.com/morikuni/failure"
)

func NewUserIDContext(userID uint) failure.Context {
	return failure.Context{"userID": fmt.Sprint(userID)}
}
