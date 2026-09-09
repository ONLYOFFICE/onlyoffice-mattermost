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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocksLoopbackWithoutAllowPrivate(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	t.Cleanup(target.Close)

	client := newHTTPClient(Options{})
	resp, err := client.Get(target.URL)
	if resp != nil {
		_ = resp.Body.Close()
	}

	require.Error(t, err)
	assert.True(t, IsBlockedDestination(err))
}

func TestAllowsPrivateWhenSet(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))

	t.Cleanup(target.Close)

	client := newHTTPClient(Options{AllowPrivate: true})
	resp, err := client.Get(target.URL)

	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRedirectHopIsRevalidated(t *testing.T) {
	private := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "private")
	}))
	t.Cleanup(private.Close)

	public := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, private.URL, http.StatusFound)
	}))
	t.Cleanup(public.Close)

	client := newHTTPClient(Options{AllowPrivate: true, MaxRedirects: 5})
	resp, err := client.Get(public.URL)

	require.NoError(t, err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	require.NoError(t, err)
	assert.Equal(t, "private", string(body))
}

func TestRedirectToBlockedDestinationRejected(t *testing.T) {
	var client *http.Client
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "ftp://example.com/x", http.StatusFound)
	}))

	t.Cleanup(first.Close)

	client = newHTTPClient(Options{AllowPrivate: true})
	resp, err := client.Get(first.URL)
	if resp != nil {
		_ = resp.Body.Close()
	}

	require.Error(t, err)
	assert.True(t, IsBlockedDestination(err))
}

func TestTooManyRedirects(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server.URL+"/next", http.StatusFound)
	}))

	t.Cleanup(server.Close)

	client := newHTTPClient(Options{AllowPrivate: true, MaxRedirects: 2})
	resp, err := client.Get(server.URL)
	if resp != nil {
		_ = resp.Body.Close()
	}

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTooManyRedirects)
}
