package util

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/oinume/lekcije/server/errors"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestEncryptString(t *testing.T) {
	a := assert.New(t)

	table := []struct {
		email string
	}{
		{"a@a.com"},
		{"oinume+test@gmail.com"},
		{"builds@travis-ci.org"},
	}
	key := strings.Repeat("a", 32)

	for _, v := range table {
		encrypted, err := EncryptString(v.email, key)
		a.Nil(err)
		a.NotEmpty(encrypted)
		fmt.Printf("%s -> %s\n", v.email, encrypted)
		//a.True(len(encrypted) > aes.BlockSize)

		decrypted, err := DecryptString(encrypted, key)
		a.Nil(err)
		a.Equal(v.email, decrypted)

		encrypted2, err := EncryptString(v.email, key)
		a.Nil(err)
		a.Equal(encrypted, encrypted2)
	}
}

func TestWriteError(t *testing.T) {
	a := assert.New(t)

	var out bytes.Buffer
	WriteError(&out, errors.Internalf("error message"))
	a.Contains(out.String(), "error message")
	a.Contains(out.String(), "github.com/oinume/lekcije/server/util.TestWriteError")
	fmt.Printf("%v\n", out.String())
}
