package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/tashfi04/printbin-server/config"
	"io"
)

func Encrypt(plainText string) (string, error) {

	c, err := aes.NewCipher([]byte(config.Crypto().Key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return hex.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string) (string, error) {

	byteCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(config.Crypto().Key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	if len(byteCipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, byteCipherText := byteCipherText[:nonceSize], byteCipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, byteCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText[:]), nil
}
