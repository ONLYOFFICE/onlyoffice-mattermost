/**
 *
 * (c) Copyright Ascensio System SIA 2023
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

func TestTimeConversion(t *testing.T) {
	assert := assert.New(t)

	converter := NewConverter()
	timestamp := converter.GetTimestamp()
	cH, cM, cS := converter.GetTime(timestamp).Clock()
	H, M, S := time.Now().Clock()

	assert.Equal(cH, H)
	assert.Equal(cM, M)
	assert.Equal(cS, S)
}
