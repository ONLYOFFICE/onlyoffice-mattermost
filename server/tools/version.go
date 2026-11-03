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
	"fmt"
	"strconv"
	"strings"
)

// DocumentServerVersion is a parsed ONLYOFFICE Docs version (e.g. "10.0.1.27" -> 10.0).
type DocumentServerVersion struct {
	Major int
	Minor int
}

// ParseDocumentServerVersion parses the version string returned by the
// Document Server "version" command. Only the numeric major and minor
// components are significant; anything after them is ignored.
func ParseDocumentServerVersion(raw string) (DocumentServerVersion, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return DocumentServerVersion{}, fmt.Errorf("empty version")
	}

	parts := strings.SplitN(raw, ".", 3)
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return DocumentServerVersion{}, fmt.Errorf("invalid major version in %q", raw)
	}

	minor := 0
	if len(parts) > 1 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil || minor < 0 {
			return DocumentServerVersion{}, fmt.Errorf("invalid minor version in %q", raw)
		}
	}

	return DocumentServerVersion{Major: major, Minor: minor}, nil
}

// AtLeast reports whether v is greater than or equal to major.minor.
func (v DocumentServerVersion) AtLeast(major, minor int) bool {
	if v.Major != major {
		return v.Major > major
	}
	return v.Minor >= minor
}
