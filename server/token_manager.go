package main

import (
	"dto"
	"fmt"

	"github.com/golang-jwt/jwt"
)

func JwtSign(config dto.Config, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, config)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func JwtDecode(jwtString string, key []byte) (string, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["field"]), nil
	} else {
		return "", err
	}
}
