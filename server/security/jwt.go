/**
 *
 * (c) Copyright Ascensio System SIA 2022
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package security

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

func JwtSign(payload jwt.Claims, key []byte) (string, error) {
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
