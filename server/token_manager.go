package main

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

type SignedObject struct {
	Field string `json:"field"`
}

func (p *Plugin) JwtSign(field string, key []byte) (string, error) {
	claims := jwt.MapClaims{
		"field": field,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss + "_" + p.configuration.DESSecret, nil
}

func (p *Plugin) JwtDecode(jwtString string, key []byte) (string, error) {
	token, err := jwt.Parse(strings.TrimSuffix(jwtString, "_"+p.configuration.DESSecret), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["field"]), nil
	} else {
		fmt.Println("INVALID TOKEN")
		return "", err
	}
}
