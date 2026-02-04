package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
)

func Encrypt(plaintextString string) (ciphertextString string, err error) {
	key := []byte(x.AESKey)
	plaintext := []byte(plaintextString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errkit.AddFuncName(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode to base64 to ensure the ciphertext is valid UTF-8 for database storage
	ciphertextString = base64.StdEncoding.EncodeToString(ciphertext)

	return ciphertextString, nil
}

func Decrypt(ciphertextString string) (plaintextString string, err error) {
	key := []byte(x.AESKey)

	// Decode from base64 first
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, data := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	plaintextString = string(plaintext)

	return plaintextString, nil
}
