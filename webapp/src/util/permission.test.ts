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

import {
    FileAccess,
    getFileAccess,
    getFilePermissions,
    getPermissionsTypeByPermissions,
} from './permission';

describe('permission utils', () => {
    it('returns available access modes', () => {
        expect(getFileAccess()).toEqual([FileAccess.EDIT_ONLY, FileAccess.READ_ONLY]);
    });

    it('maps access labels to permissions', () => {
        expect(getFilePermissions('Edit')).toEqual({edit: true});
        expect(getFilePermissions('edit')).toEqual({edit: true});
        expect(getFilePermissions('Read')).toEqual({edit: false});
        expect(getFilePermissions('unknown')).toEqual({edit: false});
    });

    it('maps permissions to access type', () => {
        expect(getPermissionsTypeByPermissions(undefined)).toBe(FileAccess.READ_ONLY);
        expect(getPermissionsTypeByPermissions({edit: true})).toBe(FileAccess.EDIT_ONLY);
        expect(getPermissionsTypeByPermissions({edit: false})).toBe(FileAccess.READ_ONLY);
    });
});
