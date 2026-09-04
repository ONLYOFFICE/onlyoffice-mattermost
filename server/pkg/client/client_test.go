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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(VersionResponse{Error: 0, Version: "8.2.0"})
	}))

	t.Cleanup(server.Close)

	c := New(crypto.NewJwtManager())
	resp, err := c.SendVersion(server.URL, VersionRequest{Command: "version"}, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Error)
	assert.Equal(t, "8.2.0", resp.Version)
}

func TestSendConvert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ConvertResponse{
			FileURL:  "https://docs.example.com/out.docx",
			FileType: "docx",
			Error:    0,
		})
	}))

	t.Cleanup(server.Close)

	c := New(crypto.NewJwtManager())
	resp, err := c.SendConvert(server.URL, ConvertRequest{
		Key:        "k1",
		Filetype:   "doc",
		Outputtype: "docx",
		URL:        "https://example.com/file.doc",
	}, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Error)
	assert.Equal(t, "docx", resp.FileType)
	assert.Contains(t, resp.FileURL, "out.docx")
}

func TestSendVersionTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_ = json.NewEncoder(w).Encode(VersionResponse{})
	}))

	t.Cleanup(server.Close)

	c := New(crypto.NewJwtManager())
	_, err := c.SendVersion(server.URL, VersionRequest{Command: "version"}, 50*time.Millisecond)
	assert.Error(t, err)
}
