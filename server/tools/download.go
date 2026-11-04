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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrForeignDownloadURL is returned when the plugin is asked to download a
// file from a host other than the configured Document Server.
var ErrForeignDownloadURL = errors.New("download url does not point to the configured document server")

// CheckDocumentServerURL verifies that target is an http(s) URL with the same
// origin as the configured Document Server address. Files the plugin downloads
// (saved documents, conversion results) are always served by the Document
// Server itself, so anything with a different origin (scheme, host or port) is
// refused. This prevents the plugin from
// being used to fetch arbitrary internal URLs (SSRF) while keeping Document
// Servers on private networks working, which is the usual self-hosted setup.
func CheckDocumentServerURL(target, documentServer string) error {
	t, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("%w: invalid url: %v", ErrForeignDownloadURL, err)
	}

	scheme := strings.ToLower(t.Scheme)
	if (scheme != "http" && scheme != "https") || t.Hostname() == "" || t.User != nil {
		return fmt.Errorf("%w: unsupported url %q", ErrForeignDownloadURL, t.Redacted())
	}

	ds, err := url.Parse(documentServer)
	if err != nil || ds.Hostname() == "" {
		return fmt.Errorf("%w: document server address is not configured", ErrForeignDownloadURL)
	}

	if origin(t) != origin(ds) {
		return fmt.Errorf("%w: url origin %q, document server origin %q", ErrForeignDownloadURL, origin(t), origin(ds))
	}

	return nil
}

// origin returns scheme://host:port with the default port made explicit, so that
// other services on the same host (different port) are not treated as the
// Document Server.
func origin(u *url.URL) string {
	scheme := strings.ToLower(u.Scheme)
	port := u.Port()
	if port == "" {
		port = "80"
		if scheme == "https" {
			port = "443"
		}
	}

	return scheme + "://" + strings.ToLower(u.Hostname()) + ":" + port
}

// NewDocumentServerClient returns an HTTP client for downloading files from the
// Document Server: redirects are followed only within the Document Server host.
// Proxy settings (HTTP_PROXY/HTTPS_PROXY/NO_PROXY) are honoured through the
// default transport.
func NewDocumentServerClient(documentServer string, timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return CheckDocumentServerURL(req.URL.String(), documentServer)
		},
	}
}
