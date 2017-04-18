package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/oinume/lekcije/server/logger"
)

func PostAPIWebhookSendGrid(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	defer r.Body.Close()
	logger.App.Info(string(body))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
