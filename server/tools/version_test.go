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

func TestParseDocumentServerVersion(t *testing.T) {
	cases := []struct {
		in      string
		major   int
		minor   int
		atLeast bool
		wantErr bool
	}{
		{"8.2.0.143", 8, 2, true, false},
		{"8.1.3", 8, 1, false, false},
		{"7.5.1.1", 7, 5, false, false},
		{"9.4.0", 9, 4, true, false},
		{"10.0.0.18", 10, 0, true, false}, // regression: "1" used to be read as the major version
		{"11", 11, 0, true, false},
		{" 12.3 ", 12, 3, true, false},
		{"", 0, 0, false, true},
		{"x.1", 0, 0, false, true},
		{"9.y", 0, 0, false, true},
	}

	for _, c := range cases {
		v, err := ParseDocumentServerVersion(c.in)
		if c.wantErr {
			assert.Error(t, err, c.in)
			continue
		}
		assert.NoError(t, err, c.in)
		assert.Equal(t, c.major, v.Major, c.in)
		assert.Equal(t, c.minor, v.Minor, c.in)
		assert.Equal(t, c.atLeast, v.AtLeast(8, 2), c.in)
	}
}
