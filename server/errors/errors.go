package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var _ = fmt.Printf

type Causer interface {
	Cause() error
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type BaseError struct {
	wrapped         error
	cause           error
	stackTrace      errors.StackTrace
	outpuStackTrace bool
}

func (be *BaseError) Cause() error {
	return be.cause
}

func (be *BaseError) StackTrace() errors.StackTrace {
	return be.stackTrace
}

func (be *BaseError) SetOutputStackTrace(output bool) {
	be.outpuStackTrace = output
}

func (be *BaseError) GetOutputStackTrace() bool {
	return be.outpuStackTrace
}

func NewBaseError(wrapped error) *BaseError {
	be := &BaseError{wrapped: wrapped}
	if c, ok := be.wrapped.(Causer); ok {
		be.cause = c.Cause()
	}
	if st, ok := be.wrapped.(StackTracer); ok {
		be.stackTrace = st.StackTrace()
	}
	return be
}

type Wrapper struct {
	*BaseError
}

func (e *Wrapper) GetOutputStackTrace() bool {
	return false
}

func (e *Wrapper) Error() string {
	return fmt.Sprintf("errors.Wrapper: %s", e.wrapped.Error())
}

func Wrapperf(err error, format string, args ...interface{}) *Wrapper {
	return &Wrapper{NewBaseError(errors.Wrapf(err, format, args...))}
}

type Internal struct {
	*BaseError
}

func (e *Internal) Error() string {
	return fmt.Sprintf("errors.Internal: %s", e.wrapped.Error())
}

type InvalidArgument struct {
	*BaseError
}

func (e *InvalidArgument) Error() string {
	return fmt.Sprintf("errors.InvalidArgument: %s", e.wrapped.Error())
}

func Internalf(format string, args ...interface{}) *Internal {
	return &Internal{NewBaseError(errors.Errorf(format, args...))}
}

func InternalWrapf(err error, format string, args ...interface{}) *Internal {
	return &Internal{NewBaseError(errors.Wrapf(err, format, args...))}
}

func InvalidArgumentf(format string, args ...interface{}) *InvalidArgument {
	return &InvalidArgument{NewBaseError(fmt.Errorf(format, args...))}
}

func InvalidArgumentWrapf(err error, format string, args ...interface{}) *InvalidArgument {
	return &InvalidArgument{NewBaseError(errors.Wrapf(err, format, args...))}
}
