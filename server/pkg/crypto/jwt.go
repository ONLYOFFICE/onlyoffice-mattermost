/**
 *
 * (c) Copyright Ascensio System SIA 2025
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
package crypto

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

type onlyofficeJwtManagerImpl struct{}

func (j onlyofficeJwtManagerImpl) Sign(key []byte, payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString(key)

	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrJwtManagerSigning, err)
	}

	return ss, nil
}

func (j onlyofficeJwtManagerImpl) Verify(key []byte, jwtToken string, body interface{}) error {
	if jwtToken == "" {
		return ErrJwtManagerEmptyToken
	}

	if body == nil {
		return ErrJwtManagerEmptyDecodingBody
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJwtManagerInvalidSigningMethod
		}

		return key, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ErrJwtManagerCastOrInvalidToken
	}

	return mapstructure.Decode(claims, body)
}
