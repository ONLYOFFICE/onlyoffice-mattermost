/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
package onlyoffice

import (
	"testing"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/api/onlyoffice/model"
	"github.com/google/uuid"
	mmModel "github.com/mattermost/mattermost-server/v6/model"
	"github.com/stretchr/testify/assert"
)

func TestGetFilePermissions(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()

	post := mmModel.Post{
		Id:      uuid.NewString(),
		UserId:  "1",
		FileIds: []string{"1"},
	}

	tests := []struct {
		name                string
		userID              string
		fileID              string
		propKey             string
		setPermissions      bool
		expectedPermissions model.Permissions
	}{
		{
			name:           "Get specific user permissions",
			propKey:        createPermissionsPropKeyName("2", "1"),
			userID:         "2",
			fileID:         "1",
			setPermissions: true,
			expectedPermissions: model.Permissions{
				Edit:    false,
				Comment: true,
			},
		}, {
			name:           "Get wildcard permissions",
			propKey:        createPermissionsPropKeyName(OnlyofficePermissionsWildcardKey, "1"),
			userID:         "2",
			fileID:         "1",
			setPermissions: true,
			expectedPermissions: model.Permissions{
				Edit:    true,
				Comment: true,
				Copy:    true,
			},
		}, {
			name:                "Get author permissions",
			propKey:             createPermissionsPropKeyName("1", "1"),
			userID:              "1",
			fileID:              "1",
			expectedPermissions: model.OnlyofficeAuthorPermissions,
		}, {
			name:                "Get user permissions",
			propKey:             createPermissionsPropKeyName("2", "1"),
			userID:              "2",
			fileID:              "1",
			expectedPermissions: model.OnlyofficeDefaultPermissions,
		},
	}

	for _, test := range tests {
		tt := test

		purgeFilePermissions(tt.fileID, &post)
		t.Run(tt.name, func(t *testing.T) {
			if tt.setPermissions {
				setFilePermission(tt.propKey, tt.expectedPermissions, &post)
			}
			assert.Equal(t, tt.expectedPermissions, helper.GetFilePermissionsByUserID(tt.userID, tt.fileID, &post))
		})
	}
}

func TestSetFilePermissions(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()
	post := mmModel.Post{
		Id:      "1",
		UserId:  "1",
		FileIds: []string{"1"},
	}

	tests := []struct {
		name                string
		userID              string
		fileID              string
		postPermissions     []model.PostPermission
		expectedPermissions model.Permissions
		size                int
	}{
		{
			name:   "Set new permissions (author)",
			userID: "1",
			fileID: "1",
			postPermissions: []model.PostPermission{
				{
					FileID:   "1",
					UserID:   "1",
					Username: "1",
					Permissions: model.Permissions{
						Edit:    true,
						Comment: true,
					},
				},
			},
			expectedPermissions: model.OnlyofficeAuthorPermissions,
			size:                0,
		}, {
			name:   "Set new permissions (user)",
			userID: "2",
			fileID: "1",
			postPermissions: []model.PostPermission{
				{
					FileID:   "1",
					UserID:   "2",
					Username: "1",
					Permissions: model.Permissions{
						Edit: true,
					},
				},
			},
			expectedPermissions: model.Permissions{
				Edit: true,
			},
			size: 1,
		}, {
			name:   "Set user permissions again",
			userID: "2",
			fileID: "1",
			postPermissions: []model.PostPermission{
				{
					FileID:   "1",
					UserID:   "2",
					Username: "1",
					Permissions: model.Permissions{
						Edit: true,
					},
				},
			},
			expectedPermissions: model.Permissions{
				Edit: true,
			},
			size: 0,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			newPermissions := helper.SetPostFilePermissions(&post, tt.postPermissions)
			assert.Equal(t, tt.size, len(newPermissions))
			assert.Equal(t, tt.expectedPermissions, helper.GetFilePermissionsByUserID(tt.userID, tt.fileID, &post))
		})
	}
}

func TestGetPostPermissionsByFileID(t *testing.T) {
	t.Parallel()

	helper := NewOnlyofficeHelper()

	getUser := func(userID string) (*mmModel.User, *mmModel.AppError) {
		return &mmModel.User{
			Id:       userID,
			Username: userID,
			Email:    userID,
		}, nil
	}

	post := mmModel.Post{
		Id:      "1",
		UserId:  "1",
		FileIds: []string{"1"},
	}

	tests := []struct {
		name           string
		userID         string
		fileID         string
		setPermissions bool
		expected       []model.UserInfoResponse
	}{
		{
			name:     "Get empty permissions",
			userID:   "1",
			fileID:   "1",
			expected: []model.UserInfoResponse{},
		}, {
			name:           "Get specific permissions",
			userID:         "2",
			fileID:         "1",
			setPermissions: true,
			expected: []model.UserInfoResponse{
				{
					ID:       "2",
					Username: "2",
					Email:    "2",
					Permissions: model.Permissions{
						Copy: true,
					},
				},
			},
		}, {
			name:           "Get wildcard permissions",
			userID:         OnlyofficePermissionsWildcardKey,
			fileID:         "1",
			setPermissions: true,
			expected: []model.UserInfoResponse{
				{
					ID:       OnlyofficePermissionsWildcardKey,
					Username: OnlyofficePermissionsWildcardKey,
					Email:    OnlyofficePermissionsWildcardKey,
					Permissions: model.Permissions{
						Print: true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		tt := test

		purgeFilePermissions(tt.fileID, &post)
		t.Run(tt.name, func(t *testing.T) {
			if tt.setPermissions {
				for _, permissions := range tt.expected {
					helper.SetPostFilePermissions(&post, []model.PostPermission{
						{
							FileID:      tt.fileID,
							UserID:      permissions.ID,
							Username:    permissions.ID,
							Permissions: permissions.Permissions,
						},
					})
				}
			}

			permissions := helper.GetPostPermissionsByFileID(tt.fileID, &post, getUser)
			assert.Equal(t, tt.expected, permissions)
		})
	}
}
