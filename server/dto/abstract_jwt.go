package dto

type JwtPayload interface {
	Valid() error
}

//TODO: Implement validation logic
func (c Config) Valid() error {
	return nil
}
func (c CommandBody) Valid() error {
	return nil
}
