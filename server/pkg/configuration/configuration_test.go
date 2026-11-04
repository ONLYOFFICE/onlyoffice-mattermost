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
package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFormats(t *testing.T) {
	formats, dropped := normalizeFormats(" DOCX, xlsx ,removedformat,, pptx ")
	assert.Equal(t, "docx,xlsx,pptx", formats)
	assert.Equal(t, []string{"removedformat"}, dropped)

	formats, dropped = normalizeFormats("gone1,gone2")
	assert.Equal(t, EmptyFormats, formats, "all dropped must mean 'none', not 'all'")
	assert.Equal(t, []string{"gone1", "gone2"}, dropped)

	formats, dropped = normalizeFormats("")
	assert.Equal(t, "", formats)
	assert.Empty(t, dropped)

	formats, dropped = normalizeFormats("none")
	assert.Equal(t, "none", formats)
	assert.Empty(t, dropped)
}

func TestSanitizeKeepsPluginActiveWithUnknownFormats(t *testing.T) {
	c := &Configuration{DESAddress: "https://docs.example.com/", DESJwt: "s", DESJwtHeader: "AuthorizationJWT", DESJwtPrefix: "Bearer", Formats: "docx,removedformat"}
	c.SanitizeConfiguration()
	assert.Equal(t, "docx", c.Formats)
	assert.Equal(t, []string{"removedformat"}, c.IgnoredFormats())
	assert.NoError(t, c.IsValid())
	assert.True(t, c.IsFormatAllowedForViewing("docx"))
	assert.False(t, c.IsFormatAllowedForViewing("xlsx"))
}
