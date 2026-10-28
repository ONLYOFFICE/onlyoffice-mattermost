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

func TestMapLanguageToTemplate(t *testing.T) {
	tests := []struct {
		locale string
		want   string
	}{
		{locale: "en-US", want: "en-US"},
		{locale: "de-DE", want: "de-DE"},
		{locale: "en", want: "en-US"},
		{locale: "ru", want: "ru-RU"},
		{locale: "zh", want: "zh-CN"},
		{locale: "unknown", want: "default"},
		{locale: "", want: "default"},
	}

	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			assert.Equal(t, tt.want, MapLanguageToTemplate(tt.locale))
		})
	}
}
