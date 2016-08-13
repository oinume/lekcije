package util

import (
	"math/rand"
	"os"
	"strconv"
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

func StringToInterfaceSlice(from ...string) []interface{} {
	to := make([]interface{}, len(from))
	for i := range from {
		to[i] = from[i]
	}
	return to
}

func StringToUint32Slice(from ...string) []uint32 {
	to := make([]uint32, len(from))
	for i := range from {
		tmp, err := strconv.ParseUint(from[i], 10, 32)
		if err != nil {
			to[i] = 0
		}
		to[i] = uint32(tmp)
	}
	return to
}

func Uint32ToInterfaceSlice(from ...uint32) []interface{} {
	to := make([]interface{}, len(from))
	for i := range from {
		to[i] = from[i]
	}
	return to
}
