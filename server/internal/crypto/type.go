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

import "github.com/golang-jwt/jwt/v5"

var _ JwtManager = (*onlyofficeJwtManager)(nil)
var _ Encoder = (*messageDigest)(nil)

type Encoder interface {
	Encode(text string) (string, error)
}

func NewMD5Encoder() Encoder {
	return messageDigest{}
}

type JwtManager interface {
	Sign(payload jwt.Claims) (string, error)
	Verify(jwt string, body interface{}) error
	GetKey() []byte
}

func NewJwtManager(key []byte) JwtManager {
	return onlyofficeJwtManager{
		key: key,
	}
}
