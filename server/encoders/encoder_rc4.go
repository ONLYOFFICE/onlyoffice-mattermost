package encoders

import (
	"crypto/rc4"
	"math/big"
)

func (e EncoderRC4) Encode(text string, key []byte) (string, error) {
	rc, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}

	textBytes := []byte(text)

	encrypted := make([]byte, len(textBytes))
	rc.XORKeyStream(encrypted, textBytes)

	data := new(big.Int).SetBytes(encrypted)

	return data.String(), nil
}

func (e EncoderRC4) Decode(text string, key []byte) (string, error) {
	sequence := new(big.Int)
	sequence.SetString(text, 10)

	decrypted := make([]byte, len(text))
	encryptedText := sequence.Bytes()

	rc, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}

	rc.XORKeyStream(decrypted, encryptedText)

	return string(decrypted), nil
}
