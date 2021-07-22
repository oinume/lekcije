package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stvp/rollbar"

	"github.com/oinume/lekcije/backend/errors"
)

/* TODO: race condition problem
==================
WARNING: DATA RACE
Read at 0x00c00014ea00 by goroutine 82:
  math/rand.(*rngSource).Uint64()
      /usr/local/go/src/math/rand/rng.go:239 +0x3e
  math/rand.(*rngSource).Int63()
      /usr/local/go/src/math/rand/rng.go:234 +0x1d9
  math/rand.(*Rand).Int63()
      /usr/local/go/src/math/rand/rand.go:85 +0x85
  math/rand.(*Rand).Int31()
      /usr/local/go/src/math/rand/rand.go:99 +0x9a
  math/rand.(*Rand).Int31n()
      /usr/local/go/src/math/rand/rand.go:134 +0x4f
  math/rand.(*Rand).Intn()
      /usr/local/go/src/math/rand/rand.go:172 +0x59
  github.com/oinume/lekcije/backend/util.RandomString()
      /go/src/github.com/oinume/lekcije/backend/util/util.go:35 +0x106
  github.com/oinume/lekcije/backend/internal/modeltest.NewUser()
      /go/src/github.com/oinume/lekcije/backend/internal/modeltest/modeltest.go:17 +0x257
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail.func1()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:117 +0x76
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail.func2()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:143 +0xd0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Previous write at 0x00c00014ea00 by goroutine 81:
  math/rand.(*rngSource).Uint64()
      /usr/local/go/src/math/rand/rng.go:239 +0x54
  math/rand.(*rngSource).Int63()
      /usr/local/go/src/math/rand/rng.go:234 +0x1d9
  math/rand.(*Rand).Int63()
      /usr/local/go/src/math/rand/rand.go:85 +0x85
  math/rand.(*Rand).Int31()
      /usr/local/go/src/math/rand/rand.go:99 +0x9a
  math/rand.(*Rand).Int31n()
      /usr/local/go/src/math/rand/rand.go:134 +0x4f
  math/rand.(*Rand).Intn()
      /usr/local/go/src/math/rand/rand.go:172 +0x59
  github.com/oinume/lekcije/backend/util.RandomString()
      /go/src/github.com/oinume/lekcije/backend/util/util.go:35 +0x106
  github.com/oinume/lekcije/backend/internal/modeltest.NewUser()
      /go/src/github.com/oinume/lekcije/backend/internal/modeltest/modeltest.go:17 +0x257
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe.func1()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:52 +0x76
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe.func2()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:78 +0xc1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Goroutine 82 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1168 +0x5bb
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:139 +0x648
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Goroutine 81 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1168 +0x5bb
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:74 +0x648
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202
==================
==================
WARNING: DATA RACE
Read at 0x00c00014ea08 by goroutine 82:
  math/rand.(*rngSource).Uint64()
      /usr/local/go/src/math/rand/rng.go:244 +0x8e
  math/rand.(*rngSource).Int63()
      /usr/local/go/src/math/rand/rng.go:234 +0x1d9
  math/rand.(*Rand).Int63()
      /usr/local/go/src/math/rand/rand.go:85 +0x85
  math/rand.(*Rand).Int31()
      /usr/local/go/src/math/rand/rand.go:99 +0x9a
  math/rand.(*Rand).Int31n()
      /usr/local/go/src/math/rand/rand.go:134 +0x4f
  math/rand.(*Rand).Intn()
      /usr/local/go/src/math/rand/rand.go:172 +0x59
  github.com/oinume/lekcije/backend/util.RandomString()
      /go/src/github.com/oinume/lekcije/backend/util/util.go:35 +0x106
  github.com/oinume/lekcije/backend/internal/modeltest.NewUser()
      /go/src/github.com/oinume/lekcije/backend/internal/modeltest/modeltest.go:17 +0x257
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail.func1()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:117 +0x76
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail.func2()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:143 +0xd0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Previous write at 0x00c00014ea08 by goroutine 81:
  math/rand.(*rngSource).Uint64()
      /usr/local/go/src/math/rand/rng.go:244 +0xaa
  math/rand.(*rngSource).Int63()
      /usr/local/go/src/math/rand/rng.go:234 +0x1d9
  math/rand.(*Rand).Int63()
      /usr/local/go/src/math/rand/rand.go:85 +0x85
  math/rand.(*Rand).Int31()
      /usr/local/go/src/math/rand/rand.go:99 +0x9a
  math/rand.(*Rand).Int31n()
      /usr/local/go/src/math/rand/rand.go:134 +0x4f
  math/rand.(*Rand).Intn()
      /usr/local/go/src/math/rand/rand.go:172 +0x59
  github.com/oinume/lekcije/backend/util.RandomString()
      /go/src/github.com/oinume/lekcije/backend/util/util.go:35 +0x106
  github.com/oinume/lekcije/backend/internal/modeltest.NewUser()
      /go/src/github.com/oinume/lekcije/backend/internal/modeltest/modeltest.go:17 +0x257
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe.func1()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:52 +0x76
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe.func2()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:78 +0xc1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Goroutine 82 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1168 +0x5bb
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_UpdateMeEmail()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:139 +0x648
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202

Goroutine 81 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1168 +0x5bb
  github.com/oinume/lekcije/backend/interface/http_test.Test_UserService_GetMe()
      /go/src/github.com/oinume/lekcije/backend/interface/http/twrip_test.go:74 +0x648
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1123 +0x202
==================
*/
var (
	letters  = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`)
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	random   = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

// TODO: Make randoms package and move to it
func RandomInt(n int) int {
	return random.Intn(n)
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
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

func Uint32ToStringSlice(from ...uint32) []string {
	to := make([]string, len(from))
	for i := range from {
		to[i] = fmt.Sprint(from[i])
	}
	return to
}

func EncryptString(plainText string, encryptionKey string) (string, error) {
	if encryptionKey == "" {
		return "", fmt.Errorf("argument encryptionKey is empty")
	}
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", errors.NewInternalError(
			errors.WithError(err),
		)
	}
	plainBytes := []byte(plainText)
	//cipherBytes := make([]byte, aes.BlockSize+len(plainBytes))
	cipherBytes := make([]byte, len(plainBytes))

	// iv = initialization vector
	iv := commonIV
	stream := cipher.NewCFBEncrypter(block, iv)
	//stream.XORKeyStream(cipherBytes[aes.BlockSize:], plainBytes)
	stream.XORKeyStream(cipherBytes, plainBytes)
	return hex.EncodeToString(cipherBytes), nil
}

func DecryptString(cipherText string, encryptionKey string) (string, error) {
	cipherBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", errors.NewInternalError(
			errors.WithError(err),
		)
	}
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", errors.NewInternalError(
			errors.WithError(err),
		)
	}
	//if len(cipherBytes) < aes.BlockSize {
	//	return "", errors.Internalf("cipherText is too short.")
	//}

	iv := commonIV
	//cipherBytes = cipherBytes[aes.BlockSize:]
	plainBytes := make([]byte, len(cipherBytes))
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plainBytes, cipherBytes)

	//plainBytes := cipherBytes
	return string(plainBytes), nil
}

func IsUserAgentPC(req *http.Request) bool {
	return !IsUserAgentSP(req) && !IsUserAgentTablet(req)
}

func IsUserAgentSP(req *http.Request) bool {
	ua := strings.ToLower(req.UserAgent())
	return strings.Contains(ua, "iphone") || strings.Contains(ua, "android") || strings.Contains(ua, "ipod")
}

func IsUserAgentTablet(req *http.Request) bool {
	ua := strings.ToLower(req.UserAgent())
	return strings.Contains(ua, "ipad")
}

func SendErrorToRollbar(err error, id string) {
	if rollbar.Token == "" {
		return
	}

	fields := make([]*rollbar.Field, 0, 10)
	if id != "" {
		fields = append(fields, &rollbar.Field{
			Name: "person",
			Data: map[string]string{
				"id": id,
			},
		})
	}

	if e, ok := err.(*errors.AnnotatedError); ok && e.OutputStackTrace() {
		stackTrace := e.StackTrace()
		frames := make([]uintptr, 0, len(stackTrace))
		for _, frame := range stackTrace {
			frames = append(frames, uintptr(frame))
		}
		stack := rollbar.BuildStackWithCallers(frames)
		rollbar.ErrorWithStack(rollbar.ERR, err, stack, fields...)
	} else {
		rollbar.Error(rollbar.ERR, err, fields...)
	}
}

func GenerateTempFileFromBase64String(dir, prefix, source string) (*os.File, error) {
	b, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return nil, errors.NewInternalError(errors.WithError(err))
	}
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, errors.NewInternalError(errors.WithError(err))
	}
	if _, err := f.Write(b); err != nil {
		return nil, errors.NewInternalError(errors.WithError(err))
	}
	return f, nil
}
