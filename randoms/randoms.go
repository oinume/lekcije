package randoms

import (
	"crypto/rand"
	"encoding/base32"
	"math/big"
)

func MustNewBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func MustNewInt64(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

func MustNewString(n int) string {
	b := MustNewBytes(n)
	return base32.StdEncoding.EncodeToString(b)
}
