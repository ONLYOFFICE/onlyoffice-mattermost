package security

import "math/big"

func encryptedTextConversion(encrypted []byte) string {
	data := new(big.Int).SetBytes(encrypted)
	return data.String()
}

func textToEncryptedConversion(text string) []byte {
	sequence := new(big.Int)
	sequence.SetString(text, 10)
	cipherBytes := sequence.Bytes()

	return cipherBytes
}
