package encoders

type Encoder interface {
	Encode(text string, key []byte) (string, error)
	Decode(text string, key []byte) (string, error)
}

//TODO: Fine tuning of individial instances (using custom constructors)
type EncoderAES struct{}
type EncoderRC4 struct{}
