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

	"github.com/oinume/lekcije/server/errors"
)

var (
	letters  = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`)
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	random   = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

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
