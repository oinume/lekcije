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

//type BaseError struct {
//	wrapped         error
//	cause           error
//	stackTrace      errors.StackTrace
//	outpuStackTrace bool
//}
//
//func (be *BaseError) Cause() error {
//	return be.cause
//}
//
//func (be *BaseError) StackTrace() errors.StackTrace {
//	return be.stackTrace
//}
//
//func (be *BaseError) SetOutputStackTrace(output bool) {
//	be.outpuStackTrace = output
//}
//
//func (be *BaseError) GetOutputStackTrace() bool {
//	return be.outpuStackTrace
//}
//
//func NewBaseError(wrapped error) *BaseError {
//	be := &BaseError{wrapped: wrapped}
//	if c, ok := be.wrapped.(Causer); ok {
//		be.cause = c.Cause()
//	}
//	if st, ok := be.wrapped.(StackTracer); ok {
//		be.stackTrace = st.StackTrace()
//	}
//	return be
//}
