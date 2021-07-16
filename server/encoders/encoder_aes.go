package encoders

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
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

	data := new(big.Int).SetBytes(encrypted)

	return data.String(), nil
}

func (e EncoderAES) Decode(text string, key []byte) (string, error) {
	sequence := new(big.Int)
	sequence.SetString(text, 10)
	cipherBytes := sequence.Bytes()

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
