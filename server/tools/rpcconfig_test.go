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

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

func TestRestoreFalseBools(t *testing.T) {
	ssl := true
	fs := model.FileSettings{AmazonS3SSL: nil, AmazonS3Trace: &ssl}

	RestoreFalseBools(&fs)
	fs.SetDefaults(true)

	assert.False(t, *fs.AmazonS3SSL, "nil after RPC means false and must not become the default true")
	assert.True(t, *fs.AmazonS3Trace, "explicit true is kept")
	assert.NotNil(t, fs.AmazonS3SignV2)

	RestoreFalseBools(nil)
	RestoreFalseBools(fs)
}
