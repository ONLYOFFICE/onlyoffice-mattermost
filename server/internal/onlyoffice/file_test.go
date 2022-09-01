/**
 *
 * (c) Copyright Ascensio System SIA 2022
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
package onlyoffice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedFileExtensions(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()

	tests := []struct {
		name      string
		ext       string
		supported bool
	}{
		{
			name:      "Supported word extension",
			ext:       "docx",
			supported: true,
		}, {
			name:      "Supported word extension (upper)",
			ext:       "Docx",
			supported: true,
		}, {
			name:      "Supported cell extension",
			ext:       "xls",
			supported: true,
		}, {
			name:      "Supported cell extension (upper)",
			ext:       "Xls",
			supported: true,
		}, {
			name:      "Supported slide extension",
			ext:       "pptx",
			supported: true,
		}, {
			name:      "Supported slide extension (upper)",
			ext:       "Pptx",
			supported: true,
		}, {
			name:      "Unsupported extension",
			ext:       "unknown",
			supported: false,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			isSupported := helper.IsExtensionSupported(tt.ext)
			assert.Equal(t, tt.supported, isSupported)
		})
	}
}

func TestEditableFileExtension(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()

	tests := []struct {
		name     string
		ext      string
		editable bool
	}{
		{
			name:     "Editable word extension",
			ext:      "docx",
			editable: true,
		}, {
			name:     "Not editable word extension",
			ext:      "Doc",
			editable: false,
		}, {
			name:     "Editable cell extension",
			ext:      "xlsx",
			editable: true,
		}, {
			name:     "Not editable cell extension",
			ext:      "Xls",
			editable: false,
		}, {
			name:     "Editable slide extension",
			ext:      "Pptx",
			editable: true,
		}, {
			name:     "Not editable slide extension",
			ext:      "pPt",
			editable: false,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			isEditable := helper.IsExtensionEditable(tt.ext)
			assert.Equal(t, tt.editable, isEditable)
		})
	}
}

func TestFileType(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()

	tests := []struct {
		name         string
		ext          string
		expectedType string
		withErr      bool
	}{
		{
			name:         "Get word type",
			ext:          "docx",
			expectedType: OnlyofficeWordType,
			withErr:      false,
		}, {
			name:         "Get cell type",
			ext:          "xls",
			expectedType: OnlyofficeCellType,
			withErr:      false,
		}, {
			name:         "Get slide type",
			ext:          "pptx",
			expectedType: OnlyofficeSlideType,
			withErr:      false,
		}, {
			name:         "Unknown type",
			ext:          "unknown",
			expectedType: "",
			withErr:      true,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			fType, err := helper.GetFileType(tt.ext)

			if tt.withErr {
				assert.ErrorIs(t, err, ErrOnlyofficeExtensionNotSupported)
			}
			assert.Equal(t, tt.expectedType, fType)
		})
	}
}
