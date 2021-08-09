import React from 'react';

import {FilePermissions, FileAccess} from 'utils/file';
import {getPermissionsMap} from 'utils/file/permissions';
import {AutocompleteUser} from 'utils/user';

export const UserRow = ({user, removeUser, changePermissions}: {user: AutocompleteUser, removeUser: (username: string) => void,
    changePermissions: (username: string, newPermission: string) => void}) => {
    const permissionsMap = getPermissionsMap();
    return (
        <div className='more-modal__row'>
            <button
                className='statuc-wrapper style--none'
                tabIndex={-1}
            >
                <span className='profile-icon'>
                    <img
                        className='Avatar Avatar-md'
                        alt={`${user.label} profile image`}
                        src={user.avatarUrl}
                    />
                </span>
            </button>
            <div
                className='more-modal__details'
                data-testid='userListItemDetails'
            >
                <div className='d-flex whitespace--nowrap'>
                    <div className='more-modal__name'>
                        <button
                            aria-label={`@${user.label}`}
                            className='user-popover style--none'
                        >
                            {user.label}
                        </button>
                    </div>
                </div>
            </div>
            <div
                data-testid='userListItemActions'
                className='more-modal__actions'
                style={{display: 'flex'}}
            >
                <select
                    value={user.permissions || FilePermissions.EDIT_ONLY.toString()}
                    onChange={(e) => changePermissions(user.label, e.target.value)}
                >
                    {permissionsMap.map((value: [FilePermissions, FileAccess]) => {
                        const permString = value[0].toString();
                        return (
                            <option
                                key={permString}
                                value={permString}
                            >
                                {permString}
                            </option>
                        );
                    })}
                </select>
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={() => removeUser(user.label)}
                >
                    <span aria-hidden='true'>{'×'}</span>
                </button>
            </div>
        </div>
    );
};
