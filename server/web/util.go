package web

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func internalServerError(w http.ResponseWriter, err error) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
	if e, ok := err.(stackTracer); ok {
		for _, f := range e.StackTrace() {
			fmt.Fprintf(w, "%+v\n", f)
		}
	}
	//}
}
