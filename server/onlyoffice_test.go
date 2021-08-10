package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"net/http/httptest"
	"security"
	"testing"
	"time"
	"utils"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func TestMockClient(t *testing.T) {
	jsonResponse := `{"success":"true"}`

	response := ioutil.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       response,
		}, nil
	}

	mockClient := MockClient{}

	apiRes, _ := mockClient.Do(&http.Request{})

	assert.NotEmpty(t, apiRes.Body)

	var decodedResp interface{}

	json.NewDecoder(apiRes.Body).Decode(&decodedResp)

	assert.Equal(t, "map[success:true]", fmt.Sprintf("%v", decodedResp))
}

func TestAuthenticationFilter(t *testing.T) {

	p, _ := initTestPlugin(t)

	authFilter := AuthenticationFilter{plugin: p}

	v := func(w http.ResponseWriter, req *http.Request) {
		cookie := http.Cookie{
			Name:  MATTERMOST_USER_COOKIE,
			Value: "userid",
		}

		req.AddCookie(&cookie)

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

func TestAesEncryptor(t *testing.T) {
	var aes security.EncryptorAES = security.EncryptorAES{}
	var encryptionKey []byte = []byte(utils.GenerateKey())
	var mockId string = "someid"

	encryptedId, _ := aes.Encrypt(mockId, encryptionKey)

	assert.NotEmpty(t, encryptedId, encryptedId)

	decryptedId, _ := aes.Decrypt(encryptedId, encryptionKey)

	assert.Equal(t, mockId, decryptedId, decryptedId)
}

func TestMD5Checksum(t *testing.T) {
	var rc4 security.EncryptorMD5 = security.EncryptorMD5{}
	var mockId string = "someid"

	encryptedId, _ := rc4.Encrypt(mockId, nil)
	encryptedIdAgain, _ := rc4.Encrypt(mockId, nil)

	assert.NotEmpty(t, encryptedId)
	assert.Equal(t, encryptedId, encryptedIdAgain)
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

	p := Plugin{}
	p.SetAPI(api)

	return &p, api
}
