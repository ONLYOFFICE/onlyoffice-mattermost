// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

import {getTranslations} from 'util/lang';
import type {FileAccess} from 'util/permission';
import {getFileAccess} from 'util/permission';
import type {MattermostUser} from 'util/user';

import React from 'react';
import Select from 'react-select';

type Props = {
    onRemoveUser: (username: string) => void;
    onChangeUserPermissions: (username: string, newPermission: string) => void;
    theme: string;
};

export const PermissionsList = (props: Props & { error: boolean; users: MattermostUser[] }) => {
    const i18n = getTranslations();
    return (
        <div
            className='more-modal__list'
            style={{padding: '0rem 1.5rem'}}
            data-theme={props.theme}
        >
            <div>
                {props.error ? (
                    <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%'}}>
                        <span style={{fontSize: '7'}}>{i18n['permissions.modal_fetch_error']}</span>
                    </div>
                ) : (
                    <>
                        {props.users.map((user) => (
                            <PermissionsRow
                                key={user.value}
                                user={user}
                                onRemoveUser={props.onRemoveUser}
                                onChangeUserPermissions={props.onChangeUserPermissions}
                                theme={props.theme}
                            />
                        ))}
                    </>
                )}
                <div id='scroller-dummy'/>
            </div>
        </div>
    );
};

const PermissionsRow = (props: Props & { user: MattermostUser }) => {
    return (
        <div
            className='more-modal__row'
            style={{padding: 0}}
            data-theme={props.theme}
        >
            <UserIcon {...props}/>
            <UserDetails {...props}/>
            <UserActions {...props}/>
        </div>
    );
};

const UserIcon = ({user}: {user: MattermostUser}) => {
    return (
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
    );
};

const UserDetails = ({user}: {user: MattermostUser}) => {
    return (
        <div className='more-modal__details'>
            <div
                className='more-modal__name'
                style={{display: 'block'}}
            >
                {`@${user.label}`}
            </div>
            <div className='more-modal__description'>
                {user.email}
            </div>
        </div>
    );
};

const UserActions = (props: Props & { user: MattermostUser }) => {
    const i18n = getTranslations();
    const permissionsMap = getFileAccess().map((entry: FileAccess) => {
        return {
            value: entry.toString(),
            label: i18n[`types.permissions.${entry.toString().toLowerCase() as 'edit' | 'read'}`] || entry.toString(),
        };
    });
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const onChange = (value: any) => {
        if (value.value) {
            props.onChangeUserPermissions(props.user.label, value.value);
        }
    };

    return (
        <div
            className='more-modal__actions'
            style={{display: 'flex', paddingRight: '0.3rem', margin: 0}}
        >
            <div style={{width: 'auto', minWidth: '113px', maxWidth: '15rem'}}>
                <Select
                    isSearchable={false}
                    value={{
                        value: props.user.fileAccess,
                        label: i18n[`types.permissions.${props.user.fileAccess.toString().toLowerCase() as 'edit' | 'read'}`] || props.user.fileAccess,
                    }}
                    options={permissionsMap}
                    onChange={onChange}
                    styles={{
                        control: (provided: any) => ({
                            ...provided,
                            backgroundColor: props.theme === 'dark' ? '#1b1d22' : provided.backgroundColor,
                            borderColor: props.theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
                        }),
                        menu: (provided: any) => ({
                            ...provided,
                            backgroundColor: props.theme === 'dark' ? '#1b1d22' : provided.backgroundColor,
                        }),
                        option: (provided: any, state: any) => ({
                            ...provided,
                            backgroundColor: props.theme === 'dark' 
                                ? state.isFocused ? 'rgba(255, 255, 255, 0.1)' : '#1b1d22'
                                : provided.backgroundColor,
                            color: props.theme === 'dark' ? '#ffffff' : provided.color,
                        }),
                        singleValue: (provided: any) => ({
                            ...provided,
                            color: props.theme === 'dark' ? '#ffffff' : provided.color,
                        }),
                    }}
                />
            </div>
            <button
                type='button'
                className='close'
                aria-label='Close'
                onClick={() => props.onRemoveUser(props.user.label)}
                style={{marginLeft: '1rem'}}
            >
                <span aria-hidden='true'>{'×'}</span>
            </button>
        </div>
    );
};
