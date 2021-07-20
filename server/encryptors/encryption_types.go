package encryptors

type Encryptor interface {
	Encrypt(text string, key []byte) (string, error)
	Decrypt(text string, key []byte) (string, error)
}

//TODO: Fine tuning of individial instances (using custom constructors)
type EncryptorAES struct{}
type EncryptorRC4 struct{}
