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
    EDIT_ONLY = 'Edit',
    READ_ONLY = 'Read'
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
