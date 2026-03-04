package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/rotisserie/eris"
)

func Aes256cbcEncode(plaintext string, key string, iv string) (string, error) {
	plaintext = string(pkcs7Pad([]byte(plaintext), aes.BlockSize))

	block, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		return "", eris.Wrap(err, "failed to create new cipher")
	}

	ciphertext := make([]byte, len(plaintext))

	for len(iv) < aes.BlockSize {
		iv += "\000"
	}

	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func Aes256cbcDecode(encodedText string, key string, iv string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", eris.Wrap(err, "failed to decode base64 string")
	}

	block, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		return "", eris.Wrap(err, "failed to create new cipher")
	}

	for len(iv) < aes.BlockSize {
		iv += "\000"
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)

	plaintext, err := pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		return "", eris.Wrap(err, "failed to unpad plaintext")
	}

	return string(plaintext), nil
}

func pkcs7Unpad(plaintext []byte, blockSize int) ([]byte, error) {
	length := len(plaintext)
	if length == 0 {
		return nil, eris.New("plaintext is empty")
	}

	if length%blockSize != 0 {
		return nil, eris.New("plaintext is not a multiple of the block size")
	}

	paddingLen := int(plaintext[length-1])
	if paddingLen > blockSize || paddingLen == 0 {
		return nil, eris.New("invalid padding size")
	}

	for _, padByte := range plaintext[length-paddingLen:] {
		if int(padByte) != paddingLen {
			return nil, eris.New("invalid padding")
		}
	}

	return plaintext[:length-paddingLen], nil
}
