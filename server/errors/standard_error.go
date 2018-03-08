package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type Code int

const (
	CodeNotFound        Code = 1
	CodeInvalidArgument Code = 2
	CodeInternal        Code = 3
)

func (c Code) String() string {
	s := "Unknown"
	switch c {
	case CodeNotFound:
		s = "NotFound"
	case CodeInvalidArgument:
		s = "InvalidArgument"
	case CodeInternal:
		s = "Internal"
	}
	return "code." + s
}

type StandardError struct {
	code             Code
	wrapped          error
	cause            error
	stackTrace       errors.StackTrace
	outputStackTrace bool
	resourceName     string
	resourceID       string
}

func NewStandardError(code Code, options ...Option) *StandardError {
	se := &StandardError{
		code:             code,
		wrapped:          errors.New(""), // As a default value
		outputStackTrace: true,
	}
	if st, ok := se.wrapped.(StackTracer); ok {
		se.stackTrace = st.StackTrace()
	}
	for _, option := range options {
		option(se)
	}
	return se
}

// Functional Option Pattern
// https://qiita.com/weloan/items/56f1c7792088b5ede136
// WithOriginalError(err), WithOutputStackTrace(false)

type Option func(*StandardError)

func WithError(err error) Option {
	return func(se *StandardError) {
		if err == nil {
			return
		}

		if st, ok := err.(StackTracer); ok {
			se.wrapped = err
			se.stackTrace = st.StackTrace()
		} else {
			// Wrap the err to save stack trace
			e := errors.WithStack(err)
			se.wrapped = err
			if st, ok := e.(StackTracer); ok {
				se.stackTrace = st.StackTrace()
			}
		}
	}
}

func WithOutputStackTrace(outputStackTrace bool) Option {
	return func(se *StandardError) {
		se.outputStackTrace = outputStackTrace
	}
}

func WithResourceName(resourceName string) Option {
	return func(se *StandardError) {
		se.resourceName = resourceName
	}
}

func WithResourceID(resourceID string) Option {
	return func(se *StandardError) {
		se.resourceID = resourceID
	}
}

func (e *StandardError) Error() string {
	return fmt.Sprintf(
		"%v: resource: name=%v, id=%v: %v",
		e.code.String(),
		e.resourceName,
		e.resourceID,
		e.wrapped.Error(),
	)
}

func (e *StandardError) Code() Code {
	return e.code
}

func (e *StandardError) StackTrace() errors.StackTrace {
	return e.stackTrace
}

func (e *StandardError) OutputStackTrace() bool {
	return e.outputStackTrace
}

func (e *StandardError) ResourceName() string {
	return e.resourceName
}

func (e *StandardError) ResourceID() string {
	return e.resourceID
}
