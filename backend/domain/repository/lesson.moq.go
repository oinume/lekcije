// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package repository

import (
	"context"
	"github.com/oinume/lekcije/backend/model2"
	"sync"
)

// Ensure, that LessonFetcherMock does implement LessonFetcher.
// If this is not the case, regenerate this file with moq.
var _ LessonFetcher = &LessonFetcherMock{}

// LessonFetcherMock is a mock implementation of LessonFetcher.
//
// 	func TestSomethingThatUsesLessonFetcher(t *testing.T) {
//
// 		// make and configure a mocked LessonFetcher
// 		mockedLessonFetcher := &LessonFetcherMock{
// 			CloseFunc: func()  {
// 				panic("mock out the Close method")
// 			},
// 			FetchFunc: func(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error) {
// 				panic("mock out the Fetch method")
// 			},
// 		}
//
// 		// use mockedLessonFetcher in code that requires LessonFetcher
// 		// and then make assertions.
//
// 	}
type LessonFetcherMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func()

	// FetchFunc mocks the Fetch method.
	FetchFunc func(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error)

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Fetch holds details about calls to the Fetch method.
		Fetch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TeacherID is the teacherID argument value.
			TeacherID uint
		}
	}
	lockClose sync.RWMutex
	lockFetch sync.RWMutex
}

// Close calls CloseFunc.
func (mock *LessonFetcherMock) Close() {
	if mock.CloseFunc == nil {
		panic("LessonFetcherMock.CloseFunc: method is nil but LessonFetcher.Close was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedLessonFetcher.CloseCalls())
func (mock *LessonFetcherMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// Fetch calls FetchFunc.
func (mock *LessonFetcherMock) Fetch(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error) {
	if mock.FetchFunc == nil {
		panic("LessonFetcherMock.FetchFunc: method is nil but LessonFetcher.Fetch was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		TeacherID uint
	}{
		Ctx:       ctx,
		TeacherID: teacherID,
	}
	mock.lockFetch.Lock()
	mock.calls.Fetch = append(mock.calls.Fetch, callInfo)
	mock.lockFetch.Unlock()
	return mock.FetchFunc(ctx, teacherID)
}

// FetchCalls gets all the calls that were made to Fetch.
// Check the length with:
//     len(mockedLessonFetcher.FetchCalls())
func (mock *LessonFetcherMock) FetchCalls() []struct {
	Ctx       context.Context
	TeacherID uint
} {
	var calls []struct {
		Ctx       context.Context
		TeacherID uint
	}
	mock.lockFetch.RLock()
	calls = mock.calls.Fetch
	mock.lockFetch.RUnlock()
	return calls
}
