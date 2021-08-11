import {FileAccess, getPermissionsTypeByPermissions} from 'utils/file';

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
