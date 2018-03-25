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

type Resource struct {
	kind  string
	key   string
	value string
}

func NewResource(kind, key, value string) *Resource {
	return &Resource{
		kind:  kind,
		key:   key,
		value: value,
	}
}

func (r *Resource) String() string {
	return fmt.Sprintf("%v:%v:%v", r.kind, r.key, r.value)
}

type AnnotatedError struct {
	code             Code
	message          string
	wrapped          error
	cause            error
	stackTrace       errors.StackTrace
	outputStackTrace bool
	resources        []*Resource
}

func NewAnnotatedError(code Code, options ...Option) *AnnotatedError {
	ae := &AnnotatedError{
		code:             code,
		wrapped:          errors.New(""), // As a default value
		outputStackTrace: true,
		resources:        make([]*Resource, 0, 20),
	}
	if st, ok := ae.wrapped.(StackTracer); ok {
		ae.stackTrace = st.StackTrace()
	}
	for _, option := range options {
		option(ae)
	}
	return ae
}

func NewInternalError(options ...Option) *AnnotatedError {
	return NewAnnotatedError(CodeInternal, options...)
}

func NewNotFoundError(options ...Option) *AnnotatedError {
	return NewAnnotatedError(CodeNotFound, options...)
}

func NewInvalidArgumentError(options ...Option) *AnnotatedError {
	return NewAnnotatedError(CodeInvalidArgument, options...)
}

// Functional Option Pattern
// https://qiita.com/weloan/items/56f1c7792088b5ede136
// WithOriginalError(err), WithOutputStackTrace(false)

type Option func(*AnnotatedError)

func WithMessage(message string) Option {
	return func(ae *AnnotatedError) {
		ae.message = message
	}
}

func WithMessagef(format string, args ...interface{}) Option {
	return func(ae *AnnotatedError) {
		ae.message = fmt.Sprintf(format, args...)
	}
}

func WithError(err error) Option {
	return func(ae *AnnotatedError) {
		if err == nil {
			return
		}

		if st, ok := err.(StackTracer); ok {
			ae.wrapped = err
			ae.stackTrace = st.StackTrace()
		} else {
			// Wrap the err to save stack trace
			e := errors.WithStack(err)
			ae.wrapped = err
			if st, ok := e.(StackTracer); ok {
				ae.stackTrace = st.StackTrace()
			}
		}
	}
}

func WithOutputStackTrace(outputStackTrace bool) Option {
	return func(ae *AnnotatedError) {
		ae.outputStackTrace = outputStackTrace
	}
}

func WithResource(r *Resource) Option {
	return func(ae *AnnotatedError) {
		ae.resources = append(ae.resources, r)
	}
}

func WithResources(resources ...*Resource) Option {
	return func(ae *AnnotatedError) {
		ae.resources = append(ae.resources, resources...)
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

func (e *AnnotatedError) Resources() []*Resource {
	return e.resources
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

func IsNotFound(err error) bool {
	if e, ok := err.(*AnnotatedError); ok {
		return e.code == CodeNotFound
	}
	return false
}
