package util

import (
	"math/rand"
	"os"
	"time"
)

var (
	letters    = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&()~|@+*[]<>/_-=^`)
	lekcijeEnv = os.Getenv("LEKCIJE_ENV")
)

func init() {
	// TODO: should use rand.New?
	rand.Seed(time.Now().UnixNano())
}

func IsProductionEnv() bool {
	return lekcijeEnv == "production"
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
