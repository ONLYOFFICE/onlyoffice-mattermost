/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testClaims struct {
	jwt.RegisteredClaims
	Key    string `json:"key" mapstructure:"key"`
	Status int    `json:"status" mapstructure:"status"`
}

func TestJwtSignAndVerifyRoundTrip(t *testing.T) {
	manager := NewJwtManager()
	key := []byte("secret-key")
	claims := testClaims{Key: "doc-key", Status: 2}

	token, err := manager.Sign(key, claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	var decoded testClaims
	err = manager.Verify(key, token, &decoded)

	require.NoError(t, err)
	assert.Equal(t, "doc-key", decoded.Key)
	assert.Equal(t, 2, decoded.Status)
}

func TestJwtVerifyEmptyToken(t *testing.T) {
	manager := NewJwtManager()
	var decoded testClaims
	err := manager.Verify([]byte("secret"), "", &decoded)
	assert.ErrorIs(t, err, ErrJwtManagerEmptyToken)
}

func TestJwtVerifyNilBody(t *testing.T) {
	manager := NewJwtManager()
	token, err := manager.Sign([]byte("secret"), testClaims{Key: "k"})
	require.NoError(t, err)

	err = manager.Verify([]byte("secret"), token, nil)
	assert.ErrorIs(t, err, ErrJwtManagerEmptyDecodingBody)
}

func TestJwtVerifyWrongKey(t *testing.T) {
	manager := NewJwtManager()
	token, err := manager.Sign([]byte("correct"), testClaims{Key: "k"})
	require.NoError(t, err)

	var decoded testClaims
	err = manager.Verify([]byte("wrong"), token, &decoded)
	assert.Error(t, err)
}

func TestJwtVerifyTamperedToken(t *testing.T) {
	manager := NewJwtManager()
	token, err := manager.Sign([]byte("secret"), testClaims{Key: "k", Status: 1})
	require.NoError(t, err)

	tampered := token[:len(token)-4] + "xxxx"
	var decoded testClaims
	err = manager.Verify([]byte("secret"), tampered, &decoded)
	assert.Error(t, err)
}
