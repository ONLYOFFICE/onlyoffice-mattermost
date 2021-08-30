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

/* eslint-disable react/jsx-no-literals */
import React from 'react';
import Select from 'react-select';

import {FilePermissions} from 'utils/file';
import {getPermissionsMap} from 'utils/file/permissions';
import {AutocompleteUser} from 'utils/user';

import {UserActions} from './permissions_user_actions';
import {UserDetails} from './permissions_user_details';
import {UserIcon} from './permissions_user_icon';

export const UserRow = ({user, removeUser, changePermissions}: {user: AutocompleteUser, removeUser: (username: string) => void,
    changePermissions: (username: string, newPermission: string) => void}) => {
    const permissionsMap = getPermissionsMap().map((entry: FilePermissions) => {
        return {
            value: entry.toString(),
            label: entry.toString(),
        };
    });
    return (
        <div className='more-modal__row'>
            <UserIcon
                alt={`${user.label} profile image`}
                src={user.avatarUrl}
            />
            <UserDetails
                username={user.label}
                email={user.email}
            />
            <UserActions>
                <div style={{width: '10rem'}}>
                    <Select
                        isSearchable={false}
                        value={{
                            value: user.permissions,
                            label: user.permissions,
                        }}
                        options={permissionsMap}
                        // eslint-disable-next-line @typescript-eslint/no-explicit-any
                        onChange={(value: any) => {
                            if (value.label) {
                                changePermissions(user.label, value.label);
                            }
                        }}
                    />
                </div>
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={() => removeUser(user.label)}
                    style={{marginLeft: '1rem'}}
                >
                    <span aria-hidden='true'>×</span>
                </button>
            </UserActions>
        </div>
    );
};
