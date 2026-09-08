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
package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeConverterGetTime(t *testing.T) {
	converter := New()
	timestamp := int64(1700000000)
	got := converter.GetTime(timestamp)

	assert.Equal(t, time.Unix(timestamp, 0), got)
}

func TestTimeConverterGetTimestamp(t *testing.T) {
	converter := New()
	before := time.Now().Unix()
	got := converter.GetTimestamp()
	after := time.Now().Unix()

	assert.GreaterOrEqual(t, got, before)
	assert.LessOrEqual(t, got, after)
}
