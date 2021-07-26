package models

type JwtPayload interface {
	Valid() error
}

func (c Config) Valid() error {
	return nil
}

func (c CommandBody) Valid() error {
	return nil
}
