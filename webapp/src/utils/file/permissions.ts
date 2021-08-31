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

import {getTranslations} from 'utils/i18n';

interface IObjectKeys {
    [key: string]: boolean | undefined;
}

export interface SubmitPermissionsPayload {
    FileId: string,
    Username: string,
    Permissions: FileAccess,
}

export interface FileAccess extends IObjectKeys {
    copy?: boolean,
    deleteCommentAuthorOnly?: boolean,
    download?: boolean,
    edit: boolean,
    editCommentAuthorOnly?: boolean,
    fillForms?: boolean,
    modifyContentControl?: boolean,
    modifyFilter?: boolean,
    print?: boolean,
    review?: boolean,
    comment?: boolean,
}

const EDIT: FileAccess = {
    edit: true,
};

const READ: FileAccess = {
    edit: false,
};

export enum FilePermissions {
    EDIT_ONLY = getTranslations()['types.permissions.edit'],
    READ_ONLY = getTranslations()['types.permissions.read']
}

const FilePermissionsMap: Map<FilePermissions, FileAccess> = new Map<FilePermissions, FileAccess>([
    [FilePermissions.EDIT_ONLY, EDIT],
    [FilePermissions.READ_ONLY, READ],
]);

export function getPermissionsMap(): FilePermissions[] {
    return [...FilePermissionsMap.keys()];
}

export function getFileAccess(permissionType: FilePermissions): FileAccess {
    return FilePermissionsMap.get(permissionType) || READ;
}

function permissionsComparison<T extends FileAccess>(firstObject: T, secondObject: T): boolean {
    const firstKeys = Object.keys(firstObject);
    const secondKeys = Object.keys(secondObject);

    if (firstKeys.length !== secondKeys.length) {
        return false;
    }

    for (const key of firstKeys) {
        const firstVal = firstObject[key];
        const secondVal = secondObject[key];
        if (firstVal !== secondVal) {
            return false;
        }
    }

    return true;
}

export function getPermissionsTypeByPermissions(permissions: FileAccess): string {
    const isEditOnly = permissionsComparison<FileAccess>(permissions, EDIT);

    return isEditOnly ? FilePermissions.EDIT_ONLY.toString() : FilePermissions.READ_ONLY.toString();
}
