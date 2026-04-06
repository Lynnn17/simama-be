package random

import (
	"crypto/rand"
	"math/big"
)

// Charset constants for code generation
const (
	AlphanumericUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphanumericLowercase = "abcdefghijklmnopqrstuvwxyz0123456789"
	AlphanumericMixed     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LettersOnly           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenerateCode generates a cryptographically secure random code of specified length using the given charset
func RandStringBytes(length int, charset string) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result), nil
}
