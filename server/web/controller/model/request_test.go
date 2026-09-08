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
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertFileRequestValidate(t *testing.T) {
	require.NoError(t, (&ConvertFileRequest{FileID: "file-1"}).Validate())
	assert.Error(t, (&ConvertFileRequest{}).Validate())
}

func TestConvertFileResponseToJSON(t *testing.T) {
	raw := (&ConvertFileResponse{Error: 0}).ToJSON()
	assert.Contains(t, string(raw), `"error":0`)
}

func TestNewFileRequestValidate(t *testing.T) {
	require.NoError(t, (&NewFileRequest{
		ChannelID: "ch-1",
		FileName:  "doc",
		FileType:  "docx",
	}).Validate())
	assert.Error(t, (&NewFileRequest{ChannelID: "ch-1"}).Validate())
}
