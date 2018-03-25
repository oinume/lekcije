package errors

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnnotatedError(t *testing.T) {
	a := assert.New(t)

	err := NewAnnotatedError(
		CodeNotFound,
		WithMessage("failed"),
		WithError(fmt.Errorf("not exist")),
		WithOutputStackTrace(false),
		WithResources(NewResource("user", "id", "12345")),
	)
	a.Equal(CodeNotFound, err.Code())
	a.Equal(err.Error(), "code.NotFound: failed: not exist")
	a.Equal("user", err.Resources()[0].kind)
	a.Equal("id", err.Resources()[0].entries[0].Key)
	a.EqualValues("12345", err.Resources()[0].entries[0].Value)

	var out bytes.Buffer
	fmt.Fprintf(&out, "%+v\n", err.StackTrace())
	a.Contains(out.String(), "github.com/oinume/lekcije/server/errors.NewAnnotatedError")
	a.Contains(out.String(), "TestNewAnnotatedError")
}
