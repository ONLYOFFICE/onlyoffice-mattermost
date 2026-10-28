// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
