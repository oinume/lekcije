package http

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oinume/lekcije/server/errors"
)

func TestInternalServerError(t *testing.T) {
	a := assert.New(t)

	w := httptest.NewRecorder()
	err := errors.NewInternalError(errors.WithError(fmt.Errorf("new error")))
	internalServerError(nil, w, err, 0)

	a.Contains(w.Body.String(), "code.Internal: new error")
}
