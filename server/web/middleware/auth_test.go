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
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/tools"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorizationMiddlewareWithHeader(t *testing.T) {
	api := &plugintest.API{}
	authMiddleware := NewAuthorizationMiddleware(api)

	called := false
	handler := authMiddleware.Handle(func(api plugin.API) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			called = true
			assert.Equal(t, "user-1", r.Header.Get(tools.MMAuthHeader))
			rw.WriteHeader(http.StatusOK)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set(tools.MMAuthHeader, "user-1")
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthorizationMiddlewareWithCode(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVGet", "auth-code").Return([]byte("user-2"), nil)
	authMiddleware := NewAuthorizationMiddleware(api)

	called := false
	handler := authMiddleware.Handle(func(api plugin.API) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			called = true
			assert.Equal(t, "user-2", r.Header.Get(tools.MMAuthHeader))
			rw.WriteHeader(http.StatusOK)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test?code=auth-code", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, recorder.Code)

	api.AssertExpectations(t)
}

func TestAuthorizationMiddlewareForbidden(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVGet", "missing").Return([]byte(""), nil)
	api.On("LogWarn", mock.Anything).Return()
	authMiddleware := NewAuthorizationMiddleware(api)

	called := false
	handler := authMiddleware.Handle(func(api plugin.API) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			called = true
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test?code=missing", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAuthorizationMiddlewareKVError(t *testing.T) {
	api := &plugintest.API{}
	api.On("KVGet", "bad").Return([]byte(nil), model.NewAppError("KVGet", "kv.get", nil, "fail", http.StatusInternalServerError))
	api.On("LogWarn", mock.Anything).Return()
	authMiddleware := NewAuthorizationMiddleware(api)

	handler := authMiddleware.Handle(func(api plugin.API) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test?code=bad", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}
