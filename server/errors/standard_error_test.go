package errors

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStandardError(t *testing.T) {
	a := assert.New(t)

	err := NewStandardError(
		CodeNotFound,
		WithError(fmt.Errorf("not exist")),
		WithOutputStackTrace(false),
		WithResourceName("user"),
		WithResourceID("12345"),
	)
	a.Equal(CodeNotFound, err.Code())
	a.Equal(err.Error(), "code.NotFound: resource: name=user, id=12345: not exist")
	a.Equal("user", err.resourceName)
	a.Equal("12345", err.ResourceID())

	var out bytes.Buffer
	fmt.Fprintf(&out, "%+v\n", err.StackTrace())
	a.Contains(out.String(), "github.com/oinume/lekcije/server/errors.NewStandardError")
	a.Contains(out.String(), "TestNewStandardError")
}
