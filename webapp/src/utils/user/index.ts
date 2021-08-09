import {FileAccess, getPermissionsTypeByPermissions} from 'utils/file';

export type User = {
    id: string,
    username: string,
    permissions: FileAccess,
};

export type AutocompleteUser = {
    value: string,
    label: string,
    avatarUrl: string,
    permissions: string,
};

const getUserAvatarUrl = (id: string): string => {
    return `/api/v4/users/${id}/image?_=0`;
};

export const mapUserToAutocompleteUser = (user: User): AutocompleteUser => {
    return {
        value: user.id,
        label: user.username,
        avatarUrl: getUserAvatarUrl(user.id),
        permissions: getPermissionsTypeByPermissions(user.permissions),
    };
};
