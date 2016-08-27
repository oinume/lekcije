package util

import (
	"crypto/aes"
	"fmt"
	"strings"
	"testing"

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
		a.NoError(err)
		a.NotEmpty(encrypted)
		a.True(len(encrypted) > aes.BlockSize)

		decrypted, err := DecryptString(encrypted, key)
		a.NoError(err)
		a.Equal(v.email, decrypted)
	}
}
