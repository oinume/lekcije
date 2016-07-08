package util

import (
	"time"
	"math/rand"
)

var (
	letters = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&()~|@+*[]<>/_-=^`)
)

func init() {
	// TODO: should use rand.New?
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
