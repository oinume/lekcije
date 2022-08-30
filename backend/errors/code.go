package errors

import (
	"github.com/morikuni/failure"
)

const (
	Internal        failure.StringCode = "Internal"
	InvalidArgument failure.StringCode = "InvalidArgument"
	NotFound        failure.StringCode = "NotFound"
	Unauthenticated failure.StringCode = "Unauthenticated"
)
