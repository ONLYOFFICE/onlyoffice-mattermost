package encoders

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func (e EncoderAES) Encode(text string, key []byte) (string, error) {
	byteText := []byte(text)

	aes, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	aesGCM, err := cipher.NewGCM(aes)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return "", err
	}
	encrypted := aesGCM.Seal(nonce, nonce, byteText, nil)

	data := encryptedTextConversion(encrypted)

	return data, nil
}

func (e EncoderAES) Decode(text string, key []byte) (string, error) {
	cipherBytes := textToEncryptedConversion(text)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", err
	}

	nonce, ciphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
