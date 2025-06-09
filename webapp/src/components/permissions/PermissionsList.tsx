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
import {getFileAccess} from 'util/permission';
import type {FileAccess} from 'util/permission';
import type {MattermostUser} from 'util/user';

import React, {useMemo} from 'react';
import Select from 'react-select';

type Props = {
    theme: string;
    darkTheme: string;
    onRemoveUser: (username: string) => void;
    onChangeUserPermissions: (username: string, newPermission: string) => void;
};

export const PermissionsList = (props: Props & { error: boolean; users: MattermostUser[] }) => {
    const i18n = getTranslations();

    return (
        <div className='more-modal__list onlyoffice-permissions__list'>
            <div>
                {props.error ? (
                    <div className='onlyoffice-permissions__error-container'>
                        <span className='onlyoffice-permissions__error-text'>{i18n['permissions.modal_fetch_error']}</span>
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
                                darkTheme={props.darkTheme}
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
        <div className='more-modal__row onlyoffice-permissions__row'>
            <UserIcon {...props}/>
            <UserDetails user={props.user}/>
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
            <div className='more-modal__name onlyoffice-permissions__name'>
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

    const styles = useMemo(() => ({
        selectStyles: {
            control: (provided: any) => {
                return {
                    ...provided,
                    minHeight: '100%',
                    backgroundColor: props.theme === 'dark' ? 'var(--center-channel-bg)' : provided.backgroundColor,
                    borderColor: props.theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
                    '&:hover': {
                        borderColor: props.theme === 'dark' ? 'rgba(255, 255, 255, 0.2)' : provided.borderColor,
                    },
                };
            },
            menu: (provided: any) => ({
                ...provided,
                backgroundColor: props.theme === 'dark' ? 'var(--center-channel-bg)' : provided.backgroundColor,
                borderColor: props.theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
                boxShadow: props.theme === 'dark' ? '0 2px 4px rgba(0, 0, 0, 0.5)' : provided.boxShadow,
            }),
            option: (provided: any, state: any) => {
                let backgroundColor = 'white';
                let color = '#3d3c40';

                if (props.theme === 'dark') {
                    backgroundColor = 'var(--center-channel-bg)';
                    color = '#ffffff';
                }

                if (state.isFocused) {
                    backgroundColor = props.theme === 'dark' ? 
                        (props.darkTheme === 'indigo' ? '#262B39' : 
                         props.darkTheme === 'onyx' ? '#2D2E33' : 'rgba(255, 255, 255, 0.1)') 
                        : '#F1F2F3';
                }

                return {
                    ...provided,
                    backgroundColor,
                    color,
                    cursor: 'pointer',
                    ':hover': {
                        backgroundColor: props.theme === 'dark' ? 
                            (props.darkTheme === 'indigo' ? '#262B39' : 
                             props.darkTheme === 'onyx' ? '#2D2E33' : 'rgba(255, 255, 255, 0.1)') 
                            : '#F1F2F3',
                        color,
                    },
                    ':active': {
                        backgroundColor: props.theme === 'dark' ? 
                            (props.darkTheme === 'indigo' ? '#262B39' : 
                             props.darkTheme === 'onyx' ? '#2D2E33' : 'rgba(255, 255, 255, 0.1)') 
                            : '#F1F2F3',
                    },
                };
            },
            singleValue: (provided: any) => ({
                ...provided,
                color: props.theme === 'dark' ? '#ffffff' : provided.color,
            }),
            input: (provided: any) => ({
                ...provided,
                color: props.theme === 'dark' ? '#ffffff' : provided.color,
            }),
        },
    }), [props.theme]);

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
        <div className='more-modal__actions onlyoffice-permissions__actions'>
            <div className='onlyoffice-permissions__select-container'>
                <Select
                    isSearchable={false}
                    value={{
                        value: props.user.fileAccess,
                        label: i18n[`types.permissions.${props.user.fileAccess.toString().toLowerCase() as 'edit' | 'read'}`] || props.user.fileAccess,
                    }}
                    options={permissionsMap}
                    onChange={onChange}
                    styles={styles.selectStyles}
                    components={{
                        IndicatorSeparator: () => null,
                    }}
                />
            </div>
            <button
                type='button'
                className='close onlyoffice-permissions__remove-button'
                aria-label='Close'
                onClick={() => props.onRemoveUser(props.user.label)}
                style={{
                    color: props.theme === 'dark' ? '#ffffff' : '#3d3c40',
                    opacity: 0.7,
                }}
            >
                <span aria-hidden='true'>{'×'}</span>
            </button>
        </div>
    );
};
