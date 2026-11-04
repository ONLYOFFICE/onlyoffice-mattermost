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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckDocumentServerURL(t *testing.T) {
	ds := "https://docs.example.com:8443"
	ok := []string{
		"https://docs.example.com:8443/cache/files/x.docx?md5=1",
		"https://DOCS.example.com:8443/cache/files/x.docx",
	}
	bad := []string{
		"http://169.254.169.254/latest/meta-data/",
		"http://127.0.0.1:8065/api/v4/users",
		"https://docs.example.com.evil.com/x",
		"https://docs.example.com:8065/api/v4/users", // another service on the same host
		"http://docs.example.com:8443/x",             // scheme downgrade
		"https://docs.example.com/x",                 // default port 443 != 8443
		"file:///etc/passwd",
		"gopher://docs.example.com/",
		"https://user:pass@docs.example.com/x",
		"",
	}
	for _, u := range ok {
		assert.NoError(t, CheckDocumentServerURL(u, ds), u)
	}
	for _, u := range bad {
		assert.ErrorIs(t, CheckDocumentServerURL(u, ds), ErrForeignDownloadURL, u)
	}
	assert.Error(t, CheckDocumentServerURL("https://docs.example.com/x", ""))
	assert.NoError(t, CheckDocumentServerURL("https://docs.example.com:443/x", "https://docs.example.com"), "explicit default port")
}

func TestDocumentServerClientBlocksForeignRedirects(t *testing.T) {
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("secret")) }))
	defer foreign.Close()

	var ds *httptest.Server
	ds = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/local":
			http.Redirect(w, r, ds.URL+"/file", http.StatusFound)
		case "/foreign":
			http.Redirect(w, r, "http://localhost:"+foreign.Listener.Addr().String()[len("127.0.0.1:"):]+"/", http.StatusFound)
		default:
			_, _ = w.Write([]byte("doc"))
		}
	}))
	defer ds.Close()

	client := NewDocumentServerClient(ds.URL, 5*time.Second)

	resp, err := client.Get(ds.URL + "/local")
	assert.NoError(t, err, "redirect within the document server is followed")
	if resp != nil {
		resp.Body.Close()
	}

	_, err = client.Get(ds.URL + "/foreign")
	assert.ErrorIs(t, err, ErrForeignDownloadURL, "redirect to another host is refused")
}
