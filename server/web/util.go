package web

import (
	"fmt"
	"net/http"

	"github.com/oinume/lekcije/server/util"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func TemplateDir() string {
	if util.IsProductionEnv() {
		return "static"
	} else {
		return "src/www"
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
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
