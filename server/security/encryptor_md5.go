package security

import (
	"crypto/md5"
	"encoding/hex"
)

func (e EncryptorMD5) Encrypt(text string, key []byte) (string, error) {
	hasher := md5.New()

	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (e EncryptorMD5) Decrypt(text string, key []byte) (string, error) {

	return text, nil
}
