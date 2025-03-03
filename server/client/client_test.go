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
package client

import (
	"testing"
)

func TestCommandClient(t *testing.T) {
	// secret := os.Getenv("JWT_SECRET")
	// invalid := "invalid"

	// tests := []struct {
	// 	name        string
	// 	command     model.CommandVersionRequest
	// 	secret      string
	// 	expectedErr int
	// }{
	// 	{
	// 		name: "Valid version command with valid secret",
	// 		command: model.CommandVersionRequest{
	// 			Command: "version",
	// 			StandardClaims: jwt.StandardClaims{
	// 				IssuedAt:  time.Now().Unix(),
	// 				ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
	// 			},
	// 		},
	// 		secret:      secret,
	// 		expectedErr: 0,
	// 	}, {
	// 		name: "Valid version command with invalid secret",
	// 		command: model.CommandVersionRequest{
	// 			Command: "version",
	// 			StandardClaims: jwt.StandardClaims{
	// 				IssuedAt:  time.Now().Unix(),
	// 				ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
	// 			},
	// 		},
	// 		secret:      invalid,
	// 		expectedErr: 6,
	// 	}, {
	// 		name: "Invalid command",
	// 		command: model.CommandVersionRequest{
	// 			Command: "invalid",
	// 			StandardClaims: jwt.StandardClaims{
	// 				IssuedAt:  time.Now().Unix(),
	// 				ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
	// 			},
	// 		},
	// 		secret:      secret,
	// 		expectedErr: 1,
	// 	},
	// }

	// for _, test := range tests {
	// 	tt := test

	// 	client := NewOnlyofficeCommandClient(false, crypto.NewJwtManager([]byte(tt.secret)))

	// 	t.Run(tt.name, func(t *testing.T) {
	// 		response, _ := client.SendVersion("https://kim.teamlab.info/coauthoring/CommandService.ashx", tt.command, 2*time.Second)
	// 		assert.Equal(t, tt.expectedErr, response.Error)

	// 		if response.Error == 0 {
	// 			assert.NotEmpty(t, response.Version)
	// 		} else {
	// 			assert.Empty(t, response.Version)
	// 		}
	// 	})
	// }
}
