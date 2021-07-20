package encryptors

import (
	"crypto/rc4"
)

func (e EncryptorRC4) Encrypt(text string, key []byte) (string, error) {
	rc, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}

	textBytes := []byte(text)

	encrypted := make([]byte, len(textBytes))
	rc.XORKeyStream(encrypted, textBytes)

	data := encryptedTextConversion(encrypted)

	return data, nil
}

func (e EncryptorRC4) Decrypt(text string, key []byte) (string, error) {
	decrypted := make([]byte, len(text))

	encryptedText := textToEncryptedConversion(text)

	rc, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}

	rc.XORKeyStream(decrypted, encryptedText)

	return string(decrypted), nil
}
