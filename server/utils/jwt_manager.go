package utils

import (
	"errors"
	"fmt"
	"models"

	"github.com/golang-jwt/jwt"
)

func JwtSign(payload models.JwtPayload, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func JwtDecode(jwtString string, key []byte) (jwt.MapClaims, error) {
	if jwtString == "" {
		return nil, errors.New("The JWT string is empty")
	}
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
