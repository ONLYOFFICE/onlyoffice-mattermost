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
package file

import (
	"encoding/json"
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/public"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/web/controller/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHelper(t *testing.T) FileHelper {
	t.Helper()
	formatManager, err := public.NewMapFormatManager()
	require.NoError(t, err)
	return New(formatManager)
}

func TestIsExtensionSupportedAndEditable(t *testing.T) {
	helper := newHelper(t)

	assert.True(t, helper.IsExtensionSupported("docx"))
	assert.True(t, helper.IsExtensionEditable("docx"))
	assert.True(t, helper.IsExtensionSupported("DOCX"))
	assert.False(t, helper.IsExtensionSupported("not-a-format"))
}

func TestGetFileType(t *testing.T) {
	helper := newHelper(t)

	fileType, err := helper.GetFileType("docx")
	require.NoError(t, err)
	assert.Equal(t, onlyofficeWordType, fileType)

	fileType, err = helper.GetFileType("xlsx")
	require.NoError(t, err)
	assert.Equal(t, onlyofficeCellType, fileType)

	_, err = helper.GetFileType("zzz")
	assert.ErrorIs(t, err, ErrExtensionNotSupported)
}

func TestGenerateKey(t *testing.T) {
	helper := newHelper(t)
	key := helper.GenerateKey()
	assert.NotEmpty(t, key)
	assert.NotContains(t, key, "-")
	assert.NotEqual(t, key, helper.GenerateKey())
}

func TestGetFilePermissionsByUserID(t *testing.T) {
	helper := newHelper(t)
	post := &mmModel.Post{UserId: "author"}
	fileID := "file-1"

	assert.Equal(t, model.OnlyofficeAuthorPermissions, helper.GetFilePermissionsByUserID("author", fileID, post))
	assert.Equal(t, model.OnlyofficeDefaultPermissions, helper.GetFilePermissionsByUserID("other", fileID, post))

	permissions := model.Permissions{Edit: true, Comment: true}
	raw, err := json.Marshal(permissions)

	require.NoError(t, err)

	post.AddProp(createPermissionsPropKeyName("other", fileID), raw)

	assert.Equal(t, permissions, helper.GetFilePermissionsByUserID("other", fileID, post))
}

func TestGetFilePermissionsWildcard(t *testing.T) {
	helper := newHelper(t)
	post := &mmModel.Post{UserId: "author"}
	fileID := "file-1"
	permissions := model.Permissions{Edit: true, Download: true}
	raw, err := json.Marshal(permissions)

	require.NoError(t, err)

	post.AddProp(createPermissionsPropKeyName(onlyofficePermissionsWildcardKey, fileID), raw)

	assert.Equal(t, permissions, helper.GetFilePermissionsByUserID("someone", fileID, post))
	assert.Equal(t, onlyofficePermissionsWildcardKey, helper.GetWildcardUser())
}

func TestSetPostFilePermissions(t *testing.T) {
	helper := newHelper(t)
	fileID := "file-1"
	post := &mmModel.Post{
		UserId:  "author",
		FileIds: mmModel.StringArray{fileID},
	}

	notify := helper.SetPostFilePermissions(post, []model.PostPermission{
		{FileID: fileID, UserID: "author", Permissions: model.OnlyofficeAuthorPermissions},
		{FileID: fileID, UserID: "editor", Permissions: model.Permissions{Edit: true}},
	})

	assert.Len(t, notify, 1)
	assert.Equal(t, "editor", notify[0].UserID)
	assert.True(t, helper.UserHasFilePermissions("editor", fileID, post))
	assert.Equal(t, model.Permissions{Edit: true}, helper.GetFilePermissionsByUserID("editor", fileID, post))
}

func TestGetPostPermissionsByFileID(t *testing.T) {
	helper := newHelper(t)
	fileID := "file-1"
	post := &mmModel.Post{UserId: "author"}
	permissions := model.Permissions{Edit: true}
	raw, err := json.Marshal(permissions)

	require.NoError(t, err)

	post.AddProp(createPermissionsPropKeyName("user-1", fileID), raw)
	post.AddProp(createPermissionsPropKeyName(onlyofficePermissionsWildcardKey, fileID), raw)

	users := helper.GetPostPermissionsByFileID(fileID, post, func(id string) (*mmModel.User, *mmModel.AppError) {
		return &mmModel.User{Id: id, Username: "u1", Email: "u1@example.com"}, nil
	})

	require.Len(t, users, 2)
	ids := map[string]bool{}
	for _, u := range users {
		ids[u.ID] = true
	}

	assert.True(t, ids["user-1"])
	assert.True(t, ids[onlyofficePermissionsWildcardKey])
}
