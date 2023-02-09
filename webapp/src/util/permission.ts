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
interface IObjectKeys {
    [key: string]: boolean | undefined;
}

export interface SubmitPermissionsRequest {
    fileID: string,
    userID: string,
    username: string,
    permissions: FilePermissions,
}

export interface FilePermissions extends IObjectKeys {
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

const EDIT: FilePermissions = {
    edit: true,
};

const READ: FilePermissions = {
    edit: false,
};

export enum FileAccess {
    EDIT_ONLY = 'Edit',
    READ_ONLY = 'Read'
}

const FilePermissionsMap: Map<FileAccess, FilePermissions> = new Map<FileAccess, FilePermissions>([
    [FileAccess.EDIT_ONLY, EDIT],
    [FileAccess.READ_ONLY, READ],
]);

export function getFileAccess(): FileAccess[] {
    return [...FilePermissionsMap.keys()];
}

export function getFilePermissions(accessType: string): FilePermissions {
    return FilePermissionsMap.get(accessType.toLowerCase() === FileAccess.EDIT_ONLY.toString().toLocaleLowerCase() ? FileAccess.EDIT_ONLY : FileAccess.READ_ONLY) || READ;
}

function permissionsComparison<T extends FilePermissions>(firstObject: T, secondObject: T): boolean {
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

export function getPermissionsTypeByPermissions(permissions: FilePermissions | undefined): FileAccess {
    if (!permissions) {
        return FileAccess.READ_ONLY;
    }
    const isEditOnly = permissionsComparison<FilePermissions>(permissions, EDIT);
    return isEditOnly ? FileAccess.EDIT_ONLY : FileAccess.READ_ONLY;
}
