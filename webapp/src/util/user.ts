/**
 *
 * (c) Copyright Ascensio System SIA 2022
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

import {ONLYOFFICE_WILDCARD_USER} from './const';
import {FileAccess, FilePermissions, getPermissionsTypeByPermissions} from './permission';

export type OnlyofficeUser = {
    id: string,
    username: string,
    permissions: FilePermissions,
    email: string,
};

export type MattermostUser = {
    value: string,
    label: string,
    avatarUrl: string,
    fileAccess: string,
    email: string,
};

export const getUserAvatarUrl = (id: string): string => {
    if (id.length < 1) {
        return '';
    }
    return `/api/v4/users/${id}/image?_=0`;
};

export const getUniqueMattermostUsers = (userProfile: UserProfile[], users: MattermostUser[]): MattermostUser[] => {
    const permissions: MattermostUser[] = [];
    userProfile.forEach((u) => {
        if (!users.find((us) => us.value === u.id)) {
            const user: MattermostUser = {
                avatarUrl: getUserAvatarUrl(u.id),
                email: u.email,
                label: u.username,
                value: u.id,
                fileAccess: FileAccess.EDIT_ONLY,
            };
            permissions.push(user);
        }
    });
    return permissions;
};

export const mapUserToMattermostUser = (user: OnlyofficeUser): MattermostUser => {
    return {
        value: user.id,
        label: user.username,
        avatarUrl: getUserAvatarUrl(user.id),
        fileAccess: getPermissionsTypeByPermissions(user.permissions),
        email: user.email,
    };
};

export const mapUsersToMattermostUsers = (users: OnlyofficeUser[]): MattermostUser[] => {
    return users.filter((user) => user.id !== ONLYOFFICE_WILDCARD_USER).map((user) => mapUserToMattermostUser(user));
};

export const sortMattermostUsers = (users: MattermostUser[]): MattermostUser[] => {
    users.sort((a, b) => {
        if (a.label < b.label) {
            return -1;
        }
        if (a.label > b.label) {
            return 1;
        }
        return 0;
    });
    return users;
};
