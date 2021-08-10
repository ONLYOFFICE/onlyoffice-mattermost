package security

type Encryptor interface {
	Encrypt(text string, key []byte) (string, error)
	Decrypt(text string, key []byte) (string, error)
}

type EncryptorAES struct{}
type EncryptorMD5 struct{}
