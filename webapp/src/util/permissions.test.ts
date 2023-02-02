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

import {FileAccess, getFileAccess, getFilePermissions, getPermissionsTypeByPermissions} from './permission';

describe('ONLYOFFICE Permissions util', () => {
    it('get supported file access', () => {
        expect(getFileAccess()).toHaveLength(2);
    });

    it('get access type by permission (edit only)', () => {
        expect(getPermissionsTypeByPermissions({
            edit: true,
        })).toBe(FileAccess.EDIT_ONLY);
    });

    it('get access type by permission (read only)', () => {
        expect(getPermissionsTypeByPermissions({
            edit: false,
        })).toBe(FileAccess.READ_ONLY);
    });

    it('get access type by permission (not supported)', () => {
        expect(getPermissionsTypeByPermissions({
            edit: true,
            copy: true,
        })).toBe(FileAccess.READ_ONLY);
    });

    it('get edit file permissions', () => {
        expect(getFilePermissions('edit')).toEqual({edit: true});
        expect(getFilePermissions('edIt')).toEqual({edit: true});
    });

    it('get read file permissions', () => {
        expect(getFilePermissions('read')).toEqual({edit: false});
        expect(getFilePermissions('ReAd')).toEqual({edit: false});
    });

    it('get read file permissions (with unknown)', () => {
        expect(getFilePermissions('unknown')).toEqual({edit: false});
    });
});
