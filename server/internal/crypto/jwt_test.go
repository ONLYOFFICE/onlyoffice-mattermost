/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtManager(t *testing.T) {
	t.Parallel()

	var manager JwtManager = NewJwtManager([]byte("secret"))

	tests := []struct {
		name          string
		command       model.CommandVersionRequest
		resultCommand string
		withErr       bool
	}{
		{
			name: "Valid command body",
			command: model.CommandVersionRequest{
				Command: "test",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
				},
			},
			withErr: false,
		}, {
			name: "Expired command body",
			command: model.CommandVersionRequest{
				Command: "bruh",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Second * 100)),
				},
			},
			withErr: true,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			dummy := model.CommandVersionRequest{}

			token, _ := manager.Sign(tt.command)

			err := manager.Verify(token, &dummy)
			if tt.withErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.command.Command, dummy.Command)
			}
		})
	}
}
