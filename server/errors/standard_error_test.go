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
		WithMessage("failed"),
		WithError(fmt.Errorf("not exist")),
		WithOutputStackTrace(false),
		WithResource("user", "id", "12345"),
	)
	a.Equal(CodeNotFound, err.Code())
	a.Equal(err.Error(), "code.NotFound: failed: not exist")
	a.Equal("user", err.ResourceKind())
	a.Equal("id", err.ResourceKey())
	a.Equal("12345", err.ResourceValue())

	var out bytes.Buffer
	fmt.Fprintf(&out, "%+v\n", err.StackTrace())
	a.Contains(out.String(), "github.com/oinume/lekcije/server/errors.NewStandardError")
	a.Contains(out.String(), "TestNewStandardError")
}
