package qrcode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/skip2/go-qrcode"
)

type QRCode struct {
	Data string
}

func (q *QRCode) Generate() ([]byte, error) {
	code, err := qrcode.New(q.Data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("error generating QR code: %w", err)
	}
	return code.PNG(256)
}

func EncryptData(data, encryptionKey string) (string, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	hash := sha256.Sum256(ciphertext)
	signedData := append(ciphertext, hash[:]...)

	return base64.StdEncoding.EncodeToString(signedData), nil
}

func DecryptAndVerifyQRCodeData(qrCodeData, encryptionKey string) (string, error) {
	signedData, err := base64.StdEncoding.DecodeString(qrCodeData)
	if err != nil {
		return "", err
	}

	ciphertext := signedData[:len(signedData)-sha256.Size]
	signature := signedData[len(signedData)-sha256.Size:]

	hash := sha256.Sum256(ciphertext)
	if !bytes.Equal(hash[:], signature) {
		return "", errors.New("invalid signature")
	}

	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
