package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oinume/lekcije/backend/errors"
)

func TestWriteError(t *testing.T) {
	a := assert.New(t)

	var out bytes.Buffer
	err := errors.NewAnnotatedError(
		errors.CodeInternal,
		errors.WithError(fmt.Errorf("error message")),
	)
	WriteError(&out, err)
	a.Contains(out.String(), "code.Internal")
	a.Contains(out.String(), "error message")
	a.Contains(out.String(), "github.com/oinume/lekcije/server/cli.TestWriteError")
	//fmt.Printf("%v\n", out.String())
}
