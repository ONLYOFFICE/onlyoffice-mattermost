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
	"github.com/stretchr/testify/require"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "https", input: "https://docs.example.com"},
		{name: "http with path", input: "http://localhost:8080/onlyoffice"},
		{name: "missing scheme", input: "docs.example.com", wantErr: true},
		{name: "empty", input: "", wantErr: true},
		{name: "relative", input: "/onlyoffice", wantErr: true},
		{name: "ftp", input: "ftp://docs.example.com", wantErr: true},
		{name: "credentials", input: "https://user:pass@docs.example.com", wantErr: true},
		{name: "query", input: "https://docs.example.com?x=1", wantErr: true},
		{name: "fragment", input: "https://docs.example.com#frag", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidURL(tt.input)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalidServerURL)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "trim and strip slashes", input: "  https://docs.example.com///  ", want: "https://docs.example.com"},
		{name: "keep path", input: "https://docs.example.com/onlyoffice/", want: "https://docs.example.com/onlyoffice"},
		{name: "strip credentials", input: "https://user:secret@docs.example.com/path", want: "https://docs.example.com/path"},
		{name: "strip query", input: "https://docs.example.com/path?foo=bar&a=1", want: "https://docs.example.com/path"},
		{name: "strip fragment", input: "https://docs.example.com/path#section", want: "https://docs.example.com/path"},
		{name: "strip all extras", input: "http://user:pass@docs.example.com:8080/oo/?q=1#f", want: "http://docs.example.com:8080/oo"},
		{name: "empty", input: "   ", want: ""},
		{name: "ftp rejected", input: "ftp://docs.example.com", wantErr: true},
		{name: "missing scheme", input: "docs.example.com", wantErr: true},
		{name: "javascript", input: "javascript:alert(1)", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeURL(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidServerURL)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			if got != "" {
				assert.NoError(t, IsValidURL(got))
			}
		})
	}
}
