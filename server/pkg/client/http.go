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
package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/doyensec/safeurl"
)

const defaultMaxRedirects = 10

var (
	ErrTooManyRedirects = errors.New("stopped after too many redirects")

	allPorts     []int
	allPortsOnce sync.Once
)

type Options struct {
	AllowPrivate bool
	Timeout      time.Duration
	MaxRedirects int
}

var privateAllowCIDRs = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"0.0.0.0/8",
	"169.254.0.0/16",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"192.88.99.0/24",
	"198.18.0.0/15",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
	"100.64.0.0/10",
	"::/128",
	"::1/128",
	"100::/64",
	"2001::/23",
	"2001:2::/48",
	"2001:db8::/32",
	"2001::/32",
	"fc00::/7",
	"fe80::/10",
	"ff00::/8",
	"2002::/16",
	"64:ff9b::/96",
	"64:ff9b:1::/48",
	"5f00::/16",
	"2001:10::/28",
	"2001:20::/28",
	"3fff::/20",
	"100:0:0:1::/64",
}

func ports() []int {
	allPortsOnce.Do(func() {
		allPorts = make([]int, 65535)
		for i := range allPorts {
			allPorts[i] = i + 1
		}
	})

	return allPorts
}

func NewHttpClient(opts Options) *http.Client {
	maxRedirects := opts.MaxRedirects
	if maxRedirects <= 0 {
		maxRedirects = defaultMaxRedirects
	}

	builder := safeurl.GetConfigBuilder().
		SetAllowedSchemes("http", "https").
		SetAllowedPorts(ports()...).
		EnableIPv6(true).
		SetCheckRedirect(func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		})

	if opts.Timeout > 0 {
		builder.SetTimeout(opts.Timeout)
	}

	if opts.AllowPrivate {
		builder.SetAllowedIPsCIDR(privateAllowCIDRs...)
	}

	wrapped := safeurl.Client(builder.Build())

	return &http.Client{
		Timeout: opts.Timeout,
		Transport: &redirectValidatingTransport{
			wrapped:      wrapped,
			maxRedirects: maxRedirects,
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

type redirectValidatingTransport struct {
	wrapped      *safeurl.WrappedClient
	maxRedirects int
}

func isRedirect(code int) bool {
	switch code {
	case http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return true
	default:
		return false
	}
}

func closeDrain(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024*1024))
	_ = resp.Body.Close()
}

func buildRedirectRequest(prev *http.Request, status int, location string) (*http.Request, error) {
	method := prev.Method
	var body io.Reader

	switch status {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther:
		method = http.MethodGet
		body = nil
	default:
		if prev.GetBody != nil {
			var err error
			body, err = prev.GetBody()
			if err != nil {
				return nil, err
			}
		}
	}

	next, err := http.NewRequestWithContext(prev.Context(), method, location, body)
	if err != nil {
		return nil, err
	}

	next.Header = prev.Header.Clone()
	return next, nil
}

func (t *redirectValidatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	current := req
	for redirects := 0; ; redirects++ {
		resp, err := t.wrapped.Do(current)
		if err != nil {
			return nil, err
		}

		if !isRedirect(resp.StatusCode) {
			return resp, nil
		}

		if redirects >= t.maxRedirects {
			closeDrain(resp)
			return nil, ErrTooManyRedirects
		}

		loc, err := resp.Location()
		closeDrain(resp)
		if err != nil {
			return nil, fmt.Errorf("redirect location: %w", err)
		}

		next, err := buildRedirectRequest(current, resp.StatusCode, loc.String())
		if err != nil {
			return nil, err
		}

		current = next
	}
}

func IsBlockedDestination(err error) bool {
	if err == nil {
		return false
	}

	var (
		allowedIP     *safeurl.AllowedIPError
		allowedHost   *safeurl.AllowedHostError
		allowedScheme *safeurl.AllowedSchemeError
		allowedPort   *safeurl.AllowedPortError
		invalidHost   *safeurl.InvalidHostError
		ipv6Blocked   *safeurl.IPv6BlockedError
		credentials   *safeurl.SendingCredentialsBlockedError
	)

	return errors.As(err, &allowedIP) ||
		errors.As(err, &allowedHost) ||
		errors.As(err, &allowedScheme) ||
		errors.As(err, &allowedPort) ||
		errors.As(err, &invalidHost) ||
		errors.As(err, &ipv6Blocked) ||
		errors.As(err, &credentials) ||
		errors.Is(err, ErrTooManyRedirects)
}
