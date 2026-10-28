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
	"errors"
	"net/url"
	"strings"
)

var ErrInvalidServerURL = errors.New("invalid server address")

func IsValidURL(address string) error {
	address = strings.TrimSpace(address)
	if address == "" {
		return ErrInvalidServerURL
	}

	u, err := url.ParseRequestURI(address)
	if err != nil {
		return ErrInvalidServerURL
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return ErrInvalidServerURL
	}

	if u.Host == "" {
		return ErrInvalidServerURL
	}

	if u.User != nil {
		return ErrInvalidServerURL
	}

	if u.RawQuery != "" || u.ForceQuery {
		return ErrInvalidServerURL
	}

	if u.Fragment != "" {
		return ErrInvalidServerURL
	}

	return nil
}

func SanitizeURL(address string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", nil
	}

	u, err := url.Parse(address)
	if err != nil {
		return "", ErrInvalidServerURL
	}

	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidServerURL
	}

	if u.Host == "" {
		return "", ErrInvalidServerURL
	}

	path := u.EscapedPath()
	for strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	cleaned := &url.URL{
		Scheme: scheme,
		Host:   u.Host,
		Path:   path,
	}

	out := cleaned.String()
	if err := IsValidURL(out); err != nil {
		return "", err
	}

	return out, nil
}
