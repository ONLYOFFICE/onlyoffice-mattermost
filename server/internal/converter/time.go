/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

import "time"

func (c converter) GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (c converter) GetTime(timestamp int64) time.Time {
	micro := timestamp / int64(time.Microsecond)
	remainder := (timestamp % int64(time.Microsecond)) * int64(time.Millisecond)
	dateTime := time.Unix(micro, remainder)

	return dateTime
}
