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
package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		isValid bool
	}{
		{
			name:    "Valid url",
			url:     "http://localhost:8650",
			isValid: true,
		},
		{
			name:    "Invalid url",
			url:     "jwt.io",
			isValid: false,
		},
		{
			name:    "Invalid suffix",
			url:     "http://localhost:8080/",
			isValid: true,
		},
		{
			name:    "Trim invalid suffix",
			url:     strings.TrimSuffix("http://localhost:8080/", "/"),
			isValid: true,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			valid := IsValidURL(tt.url)
			assert.Equal(t, tt.isValid, valid)
		})
	}
}
