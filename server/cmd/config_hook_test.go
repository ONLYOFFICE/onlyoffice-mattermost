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
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestoreMissingPluginSettings(t *testing.T) {
	current := map[string]any{"desaddress": "https://docs", "desjwt": "secret", "formats": ""}

	merged, restored := restoreMissingPluginSettings(current, nil)
	assert.Equal(t, current, merged, "whole section missing: everything restored")
	assert.Equal(t, []string{"desaddress", "desjwt", "formats"}, restored)

	merged, restored = restoreMissingPluginSettings(current, map[string]any{"desaddress": "https://new", "formats": "docx"})
	assert.Equal(t, "https://new", merged["desaddress"], "explicit values are never overridden")
	assert.Equal(t, "docx", merged["formats"])
	assert.Equal(t, "secret", merged["desjwt"])
	assert.Equal(t, []string{"desjwt"}, restored)

	merged, restored = restoreMissingPluginSettings(current, map[string]any{"desaddress": "", "desjwt": "", "formats": ""})
	assert.Empty(t, restored, "explicitly cleared values stay cleared")
	assert.Equal(t, "", merged["desjwt"])

	merged, restored = restoreMissingPluginSettings(nil, map[string]any{"a": 1})
	assert.Empty(t, restored)
	assert.Equal(t, map[string]any{"a": 1}, merged)
}
