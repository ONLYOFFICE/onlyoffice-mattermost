// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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

import {ONLYOFFICE_WILDCARD_USER} from './const';
import {FileAccess} from './permission';
import {
    getUniqueMattermostUsers,
    getUserAvatarUrl,
    mapUserToMattermostUser,
    mapUsersToMattermostUsers,
    sortMattermostUsers,
} from './user';
import type {MattermostUser, OnlyofficeUser} from './user';

describe('user utils', () => {
    it('builds avatar urls', () => {
        expect(getUserAvatarUrl('')).toBe('');
        expect(getUserAvatarUrl('abc')).toBe('/api/v4/users/abc/image?_=0');
    });

    it('filters already selected users', () => {
        const existing: MattermostUser[] = [{
            value: 'u1',
            label: 'alice',
            avatarUrl: '',
            fileAccess: FileAccess.EDIT_ONLY,
            email: 'a@ex.com',
        }];
        const result = getUniqueMattermostUsers([
            {id: 'u1', username: 'alice', email: 'a@ex.com'} as any,
            {id: 'u2', username: 'bob', email: 'b@ex.com'} as any,
        ], existing);

        expect(result).toHaveLength(1);
        expect(result[0].value).toBe('u2');
        expect(result[0].label).toBe('bob');
        expect(result[0].fileAccess).toBe(FileAccess.EDIT_ONLY);
    });

    it('maps onlyoffice users and drops wildcard', () => {
        const users: OnlyofficeUser[] = [
            {id: ONLYOFFICE_WILDCARD_USER, username: '*', email: '*', permissions: {edit: true}},
            {id: 'u1', username: 'alice', email: 'a@ex.com', permissions: {edit: false}},
        ];
        const mapped = mapUsersToMattermostUsers(users);
        expect(mapped).toHaveLength(1);
        expect(mapped[0]).toEqual(mapUserToMattermostUser(users[1]));
        expect(mapped[0].fileAccess).toBe(FileAccess.READ_ONLY);
    });

    it('sorts users by label', () => {
        const users: MattermostUser[] = [
            {value: '2', label: 'bob', avatarUrl: '', fileAccess: FileAccess.EDIT_ONLY, email: ''},
            {value: '1', label: 'alice', avatarUrl: '', fileAccess: FileAccess.EDIT_ONLY, email: ''},
            {value: '3', label: 'alice', avatarUrl: '', fileAccess: FileAccess.EDIT_ONLY, email: ''},
        ];
        expect(sortMattermostUsers(users).map((u) => u.label)).toEqual(['alice', 'alice', 'bob']);
    });
});
