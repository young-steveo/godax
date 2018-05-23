package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Signature generates a signature for requests
func Signature(message []byte, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return ``, err
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write(message)
	if err != nil {
		return ``, err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
