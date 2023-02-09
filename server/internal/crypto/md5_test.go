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
package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageDigest(t *testing.T) {
	encoder := NewMD5Encoder()
	hashOne, _ := encoder.Encode("test")
	hashTwo, _ := encoder.Encode("test")

	assert.Equal(t, hashOne, hashTwo)
}
