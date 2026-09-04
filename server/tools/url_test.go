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
package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "https", input: "https://docs.example.com", want: true},
		{name: "http with path", input: "http://localhost:8080/onlyoffice/", want: true},
		{name: "missing scheme", input: "docs.example.com", want: false},
		{name: "empty", input: "", want: false},
		{name: "relative", input: "/onlyoffice", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidURL(tt.input))
		})
	}
}
