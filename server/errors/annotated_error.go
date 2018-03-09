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

type AnnotatedError struct {
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

func NewAnnotatedError(code Code, options ...Option) *AnnotatedError {
	se := &AnnotatedError{
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

type Option func(*AnnotatedError)

func WithMessage(message string) Option {
	return func(se *AnnotatedError) {
		se.message = message
	}
}

func WithError(err error) Option {
	return func(se *AnnotatedError) {
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
	return func(se *AnnotatedError) {
		se.outputStackTrace = outputStackTrace
	}
}

func WithResource(kind, key, value string) Option {
	return func(se *AnnotatedError) {
		se.resourceKind = kind
		se.resourceKey = key
		se.resourceValue = value
	}
}

func (e *AnnotatedError) Error() string {
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

func (e *AnnotatedError) Code() Code {
	return e.code
}

func (e *AnnotatedError) StackTrace() errors.StackTrace {
	return e.stackTrace
}

func (e *AnnotatedError) OutputStackTrace() bool {
	return e.outputStackTrace
}

func (e *AnnotatedError) ResourceString() string {
	return fmt.Sprintf("%v:%v:%v", e.resourceKind, e.resourceKey, e.resourceValue)
}

func (e *AnnotatedError) ResourceKind() string {
	return e.resourceKind
}

func (e *AnnotatedError) ResourceKey() string {
	return e.resourceKey
}

func (e *AnnotatedError) ResourceValue() string {
	return e.resourceValue
}

func (e *AnnotatedError) IsNotFound() bool {
	return e.code == CodeNotFound
}

func (e *AnnotatedError) IsInternal() bool {
	return e.code == CodeInternal
}

func (e *AnnotatedError) IsInvalidArgument() bool {
	return e.code == CodeInvalidArgument
}
