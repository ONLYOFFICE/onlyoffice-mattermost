package security

import (
	"crypto/rc4"
	"errors"
)

func (e EncryptorRC4) Encrypt(text string, key []byte) (string, error) {
	rc, err := rc4.NewCipher(key)
	if err != nil {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + " Could not create a new RC4 cypher" + ONLYOFFICE_LOGGER_ENCRYPTION_SUFFIX)
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
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + " Could not create a new RC4 cypher" + ONLYOFFICE_LOGGER_DECRYPTION_SUFFIX)
	}

	rc.XORKeyStream(decrypted, encryptedText)

	return string(decrypted), nil
}
