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

func TestCallbackRequestValidate(t *testing.T) {
	valid := &CallbackRequest{
		Key:    "doc-key",
		Status: 2,
		FileID: "file-1",
	}

	require.NoError(t, valid.Validate())

	missingKey := &CallbackRequest{Status: 2, FileID: "file-1"}
	assert.Error(t, missingKey.Validate())

	missingStatus := &CallbackRequest{Key: "doc-key", FileID: "file-1"}
	assert.Error(t, missingStatus.Validate())

	missingFile := &CallbackRequest{Key: "doc-key", Status: 2}
	assert.Error(t, missingFile.Validate())
}
