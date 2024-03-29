package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/morikuni/failure"
)

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
)

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

func UintToStringSlice(from ...uint) []string { // TODO: rewrite with generics
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
		return "", failure.Wrap(err)
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
		return "", failure.Wrap(err)
	}
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", failure.Wrap(err)
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

// TODO: Move to user_agent package

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
