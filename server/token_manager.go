package main

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func (p *Plugin) JwtSign(field string, key []byte) (string, error) {
	claims := jwt.MapClaims{
		"field": field,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (p *Plugin) JwtDecode(jwtString string, key []byte) (string, error) {
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
