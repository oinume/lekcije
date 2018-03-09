package errors

import (
	"bytes"
	"fmt"
	"io"

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
	message          string
	wrapped          error
	cause            error
	stackTrace       errors.StackTrace
	outputStackTrace bool
	resourceKind     string
	resourceKey      string
	resourceValue    string
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

func WithMessage(message string) Option {
	return func(se *StandardError) {
		se.message = message
	}
}

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

func WithResource(kind, key, value string) Option {
	return func(se *StandardError) {
		se.resourceKind = kind
		se.resourceKey = key
		se.resourceValue = value
	}
}

func (e *StandardError) Error() string {
	var b bytes.Buffer
	io.WriteString(&b, e.code.String())
	if e.message != "" {
		fmt.Fprintf(&b, ": %v", e.message)
	}
	if e.wrapped != nil {
		fmt.Fprintf(&b, ": %v", e.wrapped.Error())
	}
	return b.String()
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

func (e *StandardError) ResourceString() string {
	return fmt.Sprintf("%v:%v:%v", e.resourceKind, e.resourceKey, e.resourceValue)
}

func (e *StandardError) ResourceKind() string {
	return e.resourceKind
}

func (e *StandardError) ResourceKey() string {
	return e.resourceKey
}

func (e *StandardError) ResourceValue() string {
	return e.resourceValue
}

func (e *StandardError) IsNotFound() bool {
	return e.code == CodeNotFound
}

func (e *StandardError) IsInternal() bool {
	return e.code == CodeInternal
}

func (e *StandardError) IsInvalidArgument() bool {
	return e.code == CodeInvalidArgument
}
