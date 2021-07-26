package security

import (
	"errors"
	"models"

	"github.com/golang-jwt/jwt"
)

func JwtSign(payload models.JwtPayload, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString(key)

	if err != nil {
		return "", errors.New(ONLYOFFICE_LOGGER_PREFIX + "JWT Could not create a signed string with the given key")
	}
	return ss, nil
}

func JwtDecode(jwtString string, key []byte) (jwt.MapClaims, error) {
	if jwtString == "" {
		return nil, errors.New(ONLYOFFICE_LOGGER_PREFIX + "JWT string is empty")
	}

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ONLYOFFICE_LOGGER_PREFIX + "Unexpected JWT signing method")
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(ONLYOFFICE_LOGGER_PREFIX + "JWT token is not valid")
	}
}
