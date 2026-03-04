package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func EncodeHmacSha256(
	secret string,
	publicKey string,
) string {
	digest := hmac.New(sha256.New, []byte(secret))
	digest.Write([]byte(publicKey))

	return hex.EncodeToString(digest.Sum(nil))
}

func ValidateSignatureHmacSha256(
	signature string,
	publicKey string,
	secret string,
) bool {
	hashToCompute := EncodeHmacSha256(secret, publicKey)

	return hmac.Equal([]byte(signature), []byte(hashToCompute))
}
