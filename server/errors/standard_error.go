package errors

import "github.com/pkg/errors"

type Code int

const (
	CodeNotFound        Code = 1
	CodeInvalidArgument Code = 2
	CodeInternal        Code = 3
)

type StandardError struct {
	code             Code
	wrapped          error
	cause            error
	stackTrace       errors.StackTrace
	outputStackTrace bool
}

func WrapStandardError(code Code, original error) *StandardError {
	be := &StandardError{
		code:    code,
		wrapped: original,
	}
	if c, ok := be.wrapped.(Causer); ok {
		be.cause = c.Cause()
	}
	if st, ok := be.wrapped.(StackTracer); ok {
		be.stackTrace = st.StackTrace()
	}
	return be
}

// Functional Option Pattern
// https://qiita.com/weloan/items/56f1c7792088b5ede136
// WithOriginalError(err), WithOutputStackTrace(false)
