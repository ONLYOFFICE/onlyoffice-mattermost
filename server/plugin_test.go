/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/utils"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/security"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/models"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationFilter(t *testing.T) {

	p, _ := initTestPlugin(t)

	authFilter := AuthenticationFilter{plugin: p}

	v := func(w http.ResponseWriter, req *http.Request) {
		req.Header.Add(MATTERMOST_USER_HEADER, "userid")

		authFilter.DoFilter(w, req)

		fmt.Fprintln(w, req.Header.Get(ONLYOFFICE_AUTHORIZATION_USERID_HEADER))
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	header, _ := io.ReadAll(res.Body)
	res.Body.Close()

	assert.NotEmpty(t, string(header), string(header))
}

func TestFileValidationFilter(t *testing.T) {
	p, _ := initTestPlugin(t)

	fileFilter := FileValidationFilter{plugin: p}

	v := func(w http.ResponseWriter, req *http.Request) {
		cookie := http.Cookie{
			Name:  ONLYOFFICE_FILEVALIDATION_FILEID_HEADER,
			Value: "fileid",
		}

		req.AddCookie(&cookie)

		fileFilter.DoFilter(w, req)

		fmt.Fprintln(w, req.Header.Get(ONLYOFFICE_FILEVALIDATION_POSTID_HEADER))
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	header, _ := io.ReadAll(res.Body)
	res.Body.Close()

	assert.NotEmpty(t, string(header), string(header))
	assert.Equal(t, "1\n", string(header))
}

func TestPostAuthorizationFilter(t *testing.T) {
	p, _ := initTestPlugin(t)

	postAuthorizationFilter := PostAuthorizationFilter{plugin: p}

	v := func(w http.ResponseWriter, req *http.Request) {
		userCookie := http.Cookie{
			Name:  ONLYOFFICE_AUTHORIZATION_USERID_HEADER,
			Value: "userid",
		}
		postCookie := http.Cookie{
			Name:  ONLYOFFICE_FILEVALIDATION_POSTID_HEADER,
			Value: "postid",
		}

		req.AddCookie(&userCookie)
		req.AddCookie(&postCookie)

		postAuthorizationFilter.DoFilter(w, req)

		fmt.Fprintln(w, req.Header.Get(ONYLOFFICE_CHANNELVALIDATION_CHANNELID_HEADER))
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	header, _ := io.ReadAll(res.Body)
	res.Body.Close()

	assert.NotEmpty(t, string(header), string(header))
	assert.Equal(t, "Mock\n", string(header))
}

func TestBodyJwtFilter(t *testing.T) {

	p, _ := initTestPlugin(t)

	bodyJwtFilter := BodyJwtFilter{plugin: p}

	type TokenBody struct {
		Token string `json:"token,omitempty"`
	}

	v := func(w http.ResponseWriter, req *http.Request) {
		mockJwt := testMockJwt{
			Payload: "mock",
		}

		token, _ := security.JwtSign(mockJwt, []byte(p.configuration.DESJwt))

		payload := TokenBody{
			Token: token,
		}

		payloadBytes, _ := json.Marshal(payload)

		r := io.NopCloser(strings.NewReader(string(payloadBytes)))

		req.Body = r

		bodyJwtFilter.DoFilter(w, req)

		fmt.Fprintln(w, bodyJwtFilter.hasError)
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	body, _ := io.ReadAll(res.Body)
	isError, _ := strconv.ParseBool(string(body))
	res.Body.Close()

	assert.False(t, isError, isError)
}

func TestHeaderJwtFilter(t *testing.T) {
	p, _ := initTestPlugin(t)

	headerJwtFilter := HeaderJwtFilter{plugin: p}

	v := func(w http.ResponseWriter, req *http.Request) {
		mockJwt := testMockJwt{
			Payload: "mock",
		}

		token, _ := security.JwtSign(mockJwt, []byte(p.configuration.DESJwt))

		req.Header.Set(p.configuration.DESJwtHeader, p.configuration.DESJwtPrefix+token)

		headerJwtFilter.DoFilter(w, req)

		fmt.Fprintln(w, headerJwtFilter.hasError)
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	body, _ := io.ReadAll(res.Body)
	isError, _ := strconv.ParseBool(string(body))
	res.Body.Close()

	assert.False(t, isError, isError)
}

func TestUserAccessMiddleware(t *testing.T) {
	p, _ := initTestPlugin(t)

	authentication := &AuthenticationFilter{plugin: p}
	checkFile := &FileValidationFilter{plugin: p}
	postAccess := &PostAuthorizationFilter{plugin: p}

	authentication.SetNext(checkFile).SetNext(postAccess)

	v := func(w http.ResponseWriter, req *http.Request) {
		req.Header.Add(MATTERMOST_USER_HEADER, "userid")

		authentication.DoFilter(w, req)

		fmt.Fprintln(w, authentication.HasError())
	}

	ts := httptest.NewServer(http.HandlerFunc(v))
	client := ts.Client()
	res, _ := client.Get(ts.URL)

	body, _ := io.ReadAll(res.Body)
	isError, _ := strconv.ParseBool(string(body))
	res.Body.Close()

	assert.False(t, isError, isError)
}

func TestTimeConversion(t *testing.T) {
	assert := assert.New(t)

	timestamp := utils.GetTimestamp()
	cH, cM, cS := utils.GetTime(timestamp).Clock()
	H, M, S := time.Now().Clock()

	assert.Equal(cH, H)
	assert.Equal(cM, M)
	assert.Equal(cS, S)
}

func TestPermissionsConversions(t *testing.T) {
	mock := models.Permissions{
		Edit: true,
	}
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	json.NewEncoder(encoder).Encode(mock)
	encoder.Close()

	assert.NotEmpty(t, buf.String())

	decoded, _ := utils.ConvertBase64ToPermissions(buf.String())

	assert.Equal(t, mock.Edit, decoded.Edit)
}

func TestMD5Checksum(t *testing.T) {
	var rc4 security.EncryptorMD5 = security.EncryptorMD5{}
	var mockId string = "someid"

	encryptedId, _ := rc4.Encrypt(mockId, nil)
	encryptedIdAgain, _ := rc4.Encrypt(mockId, nil)

	assert.NotEmpty(t, encryptedId)
	assert.Equal(t, encryptedId, encryptedIdAgain)
}

type testMockJwt struct {
	Payload string
}

func (mock testMockJwt) Valid() error {
	return nil
}

func TestJwt(t *testing.T) {
	jwt := testMockJwt{
		Payload: "mock",
	}

	key := []byte(utils.GenerateKey())

	jwtEncoded, _ := security.JwtSign(jwt, key)

	assert.NotEmpty(t, jwtEncoded, jwtEncoded)

	jwtDecoded, _ := security.JwtDecode(jwtEncoded, key)

	body := testMockJwt{}

	mapstructure.Decode(jwtDecoded, &body)

	assert.Equal(t, jwt.Payload, body.Payload)
}

func TestGetFilePermissionsByUser(t *testing.T) {
	post := model.Post{
		UserId: "userid",
	}

	permissions, err := GetFilePermissionsByUser("userid", "fileid", post)

	assert.Empty(t, err)

	assert.Equal(t, models.ONLYOFFICE_AUTHOR_PERMISSIONS, permissions, permissions)

	propName := utils.CreateUserPermissionsPropName("fileid", "userid")

	post.AddProp(propName, permissions)

	userPermissions, err := GetFilePermissionsByUser("userid", "fileid", post)

	assert.Empty(t, err)

	assert.Equal(t, permissions, userPermissions, userPermissions)
}

func TestSetFilePermissions(t *testing.T) {
	post := model.Post{
		UserId: "userid",
	}

	permissions := models.Permissions{
		Edit: true,
	}

	SetFilePermissions(&post, "mockito", permissions)

	userPermissions, err := GetFilePermissionsByUser("userid", "fileid", post)

	assert.Empty(t, err)

	assert.Equal(t, permissions.Edit, userPermissions.Edit, userPermissions)
}

func TestFileInfos(t *testing.T) {
	_, api := initTestPlugin(t)

	fileInfos, _ := api.GetFileInfos(0, 0, nil)

	assert.Equal(t, 3, len(fileInfos))
}

func TestUsers(t *testing.T) {
	_, api := initTestPlugin(t)

	users, _ := api.GetUsers(nil)

	assert.Equal(t, 2, len(users))
}

func initTestPlugin(t *testing.T) (*Plugin, *plugintest.API) {
	api := &plugintest.API{}
	api.On("RegisterCommand", mock.Anything).Return(nil)
	api.On("UnregisterCommand", mock.Anything, mock.Anything).Return(nil)
	api.On("GetUser", mock.Anything).Return(&model.User{
		Id:       "userid",
		Nickname: "User",
	}, (*model.AppError)(nil))
	api.On("GetUsers", mock.Anything).Return([]*model.User{
		{
			Id:       "userid",
			Nickname: "User",
		},
		{
			Id:       "assistantid",
			Nickname: "Assistant",
		},
	}, (*model.AppError)(nil))

	api.On("GetFileInfos", mock.Anything, mock.Anything, mock.Anything).Return([]*model.FileInfo{
		{
			Id:        "1",
			PostId:    "1",
			Name:      "test.docx",
			Path:      "mock/files/test.docx",
			ChannelId: "1",
			Size:      1,
			Extension: "docx",
		},
		{
			Id:        "2",
			PostId:    "1",
			Name:      "test1",
			Path:      "mock/files/test1.docx",
			ChannelId: "1",
			Size:      1,
			Extension: "docx",
		},
		{
			Id:        "3",
			PostId:    "2",
			Name:      "test2",
			Path:      "mock/files/test2.docx",
			ChannelId: "1",
			Size:      1,
			Extension: "docx",
		},
	}, (*model.AppError)(nil))

	api.On("GetFileInfo", mock.Anything).Return(&model.FileInfo{
		Id:        "1",
		PostId:    "1",
		Extension: "docx",
	}, (*model.AppError)(nil))

	api.On("GetPost", mock.Anything).Return(&model.Post{
		Id:        "postid",
		UserId:    "userid",
		Message:   "Mock message",
		ChannelId: "Mock",
	}, (*model.AppError)(nil))

	api.On("GetChannelMember", mock.Anything, mock.Anything).Return(&model.ChannelMember{
		ChannelId: "Mock",
	}, (*model.AppError)(nil))

	p := Plugin{}
	p.SetAPI(api)
	p.configuration = &configuration{}
	p.configuration.DESJwt = "Mockerito"
	p.configuration.DESJwtHeader = "Mock"
	p.configuration.DESJwtPrefix = "Bearer "

	return &p, api
}
