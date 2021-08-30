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

import {UserProfile} from 'mattermost-redux/types/users';

import {FileAccess, FilePermissions, getPermissionsTypeByPermissions} from 'utils/file';

export type User = {
    id: string,
    username: string,
    permissions: FileAccess,
    email: string,
};

export type AutocompleteUser = {
    value: string,
    label: string,
    avatarUrl: string,
    permissions: string,
    email: string,
};

const getUserAvatarUrl = (id: string): string => {
    return `/api/v4/users/${id}/image?_=0`;
};

export const getUniqueAutocompleteUsers = (userProfile: UserProfile[], users: AutocompleteUser[]): AutocompleteUser[] => {
    const permissions: AutocompleteUser[] = [];
    userProfile.forEach((u) => {
        if (!users.find((us) => us.value === u.id)) {
            const user: AutocompleteUser = {
                avatarUrl: getUserAvatarUrl(u.id),
                email: u.email,
                label: u.username,
                value: u.id,
                permissions: FilePermissions.EDIT_ONLY.toString(),
            };
            permissions.push(user);
        }
    });
    return permissions;
};

export const mapUserToAutocompleteUser = (user: User): AutocompleteUser => {
    return {
        value: user.id,
        label: user.username,
        avatarUrl: getUserAvatarUrl(user.id),
        permissions: getPermissionsTypeByPermissions(user.permissions),
        email: user.email,
    };
};

export const sortAutocompleteUsers = (users: AutocompleteUser[]) => {
    users.sort((a, b) => {
        if (a.label < b.label) {
            return -1;
        }
        if (a.label > b.label) {
            return 1;
        }
        return 0;
    });
};
