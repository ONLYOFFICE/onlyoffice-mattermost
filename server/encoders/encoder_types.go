package encoders

type Encoder interface {
	Encode(text string, key []byte) (string, error)
	Decode(text string, key []byte) (string, error)
}

//TODO: Fine tuning of individial instances
type EncoderAES struct{}
type EncoderMD5 struct{}
