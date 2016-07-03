package web

import (
	"net/http"
)

func internalServerError(w http.ResponseWriter, error string) {
	http.Error(w, error, http.StatusInternalServerError)
}
