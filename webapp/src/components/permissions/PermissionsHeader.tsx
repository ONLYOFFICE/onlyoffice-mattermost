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

import {debounceUsersLoad} from 'util/func';
import {getTranslations} from 'util/lang';
import {getFileAccess} from 'util/permission';
import type {FileAccess} from 'util/permission';
import type {MattermostUser} from 'util/user';

import React, {useState, useEffect, useMemo} from 'react';
import {Button} from 'react-bootstrap';
import Select from 'react-select';
import AsyncSelect from 'react-select/async';

import type {Channel} from 'mattermost-redux/types/channels';
import type {FileInfo} from 'mattermost-redux/types/files';

type Props = {
    loading: boolean;
    channel: Channel | undefined | null;
    fileInfo: FileInfo;
    wildcardAccess: string;
    users: MattermostUser[];
    onSetWildcardAccess: (value: any) => void;
    onAppendUsers: (newUsers: MattermostUser[]) => void;
    theme: string;
    darkTheme: string | undefined;
};

export const PermissionsHeader: React.FC<Props> = ({
    loading,
    channel,
    fileInfo,
    wildcardAccess,
    users,
    onSetWildcardAccess,
    onAppendUsers,
    theme,
    darkTheme,
}) => {
    const i18n = getTranslations();
    const styles = useMemo(() => ({
        selectStyles: {
            container: (provided: any) => ({...provided, height: '100%'}),
            control: (provided: any) => ({
                ...provided,
                minHeight: '100%',
                backgroundColor: theme === 'dark' ? (darkTheme === 'indigo' ? '#1b1d22' : '#23272f') : provided.backgroundColor,
                borderColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
                '&:hover': {
                    borderColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.2)' : provided.borderColor,
                },
            }),
            menu: (provided: any) => ({
                ...provided,
                backgroundColor: theme === 'dark' ? (darkTheme === 'indigo' ? '#1b1d22' : '#23272f') : provided.backgroundColor,
                borderColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
                boxShadow: theme === 'dark' ? '0 2px 4px rgba(0, 0, 0, 0.5)' : provided.boxShadow,
            }),
            option: (provided: any, state: any) => {
                let backgroundColor;
                if (theme === 'dark') {
                    backgroundColor = state.isFocused ? 'rgba(255, 255, 255, 0.1)' : (darkTheme === 'indigo' ? '#1b1d22' : '#23272f');
                } else {
                    backgroundColor = provided.backgroundColor;
                }

                return {
                    ...provided,
                    backgroundColor,
                    color: theme === 'dark' ? '#ffffff' : provided.color,
                    '&:hover': {
                        backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.15)' : provided.backgroundColor,
                    },
                };
            },
            singleValue: (provided: any) => ({
                ...provided,
                color: theme === 'dark' ? '#ffffff' : provided.color,
            }),
            input: (provided: any) => ({
                ...provided,
                color: theme === 'dark' ? '#ffffff' : provided.color,
            }),
            placeholder: (provided: any) => ({
                ...provided,
                color: theme === 'dark' ? 'rgba(255, 255, 255, 0.5)' : provided.color,
            }),
            multiValue: (provided: any) => ({
                ...provided,
                backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : '#f0f0f0',
                borderRadius: '49px',
                margin: '2px 4px',
                padding: '2px 4px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: '6px',
            }),
            multiValueLabel: (provided: any) => ({
                ...provided,
                textAlign: 'center',
                color: theme === 'dark' ? '#ffffff' : '#3d3c40',
                fontSize: '12px',
                fontWeight: 400,
                lineHeight: '16px',
                padding: 0,
            }),
            multiValueRemove: (provided: any) => ({
                ...provided,
                width: '10px',
                height: '10px',
                minWidth: '10px',
                minHeight: '10px',
                borderRadius: '50%',
                margin: 0,
                padding: 0,
                fontSize: '0.8rem',
                border: theme === 'dark' ? '1px solid rgba(255, 255, 255, 0.3)' : '1px solid #ababad',
                backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.2)' : '#ababad',
                color: theme === 'dark' ? '#ffffff' : '#f0f0f0',
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                ':hover': {
                    backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.3)' : '#9c9c9e',
                    border: theme === 'dark' ? '1px solid rgba(255, 255, 255, 0.4)' : '1px solid #9c9c9e',
                },
            }),
        },
        permissionSelectStyles: {
            control: (provided: any) => ({
                ...provided,
                width: 'auto',
                height: '32px',
                border: 'none',
                borderRadius: '4px',
                boxShadow: 'none',
                backgroundColor: theme === 'dark' ? 'var(--center-channel-bg)' : 'transparent',
                cursor: 'pointer',
                display: 'flex',
                justifyContent: 'flex-end',
                padding: '4px 10px 5px 12px',
                ':hover': {
                    backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : '#1C58D914',
                },
            }),
            valueContainer: (provided: any) => ({
                ...provided,
                padding: 0,
                display: 'flex',
                justifyContent: 'flex-end',
            }),
            indicatorsContainer: (provided: any) => ({
                ...provided,
                padding: 0,
                display: 'flex',
                justifyContent: 'flex-end',
            }),
            singleValue: (provided: any) => ({
                ...provided,
                color: theme === 'dark' ? '#ffffff' : '#1C58D9',
                marginRight: '6px',
            }),
            dropdownIndicator: (provided: any) => ({
                ...provided,
                color: theme === 'dark' ? '#ffffff' : '#1C58D9',
                padding: 0,
                marginRight: '0px',
                ':hover': {
                    color: theme === 'dark' ? '#ffffff' : '#1C58D9',
                },
                svg: {
                    width: '14px',
                    height: '14px',
                    fill: theme === 'dark' ? '#ffffff' : '#1C58D9',
                    ':hover': {
                        fill: theme === 'dark' ? '#ffffff' : '#1C58D9',
                    },
                },
            }),
            menu: (provided: any) => ({
                ...provided,
                backgroundColor: theme === 'dark' ? 'var(--center-channel-bg)' : provided.backgroundColor,
                borderColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : provided.borderColor,
            }),
            option: (provided: any, state: any) => {
                return {
                    ...provided,
                    backgroundColor: theme === 'dark' 
                        ? (state.isFocused ? 'rgba(255, 255, 255, 0.1)' : 'var(--center-channel-bg)')
                        : (state.isSelected ? '#1C58D9' : state.isFocused ? '#1C58D914' : provided.backgroundColor),
                    color: theme === 'dark' ? '#ffffff' : (state.isSelected ? '#ffffff' : '#1C58D9'),
                    ':hover': {
                        backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : (state.isSelected ? '#1C58D9' : '#1C58D914'),
                        color: theme === 'dark' ? '#ffffff' : (state.isSelected ? '#ffffff' : '#1C58D9'),
                    },
                    ':active': {
                        backgroundColor: theme === 'dark' ? 'rgba(255, 255, 255, 0.1)' : (state.isSelected ? '#1C58D9' : '#1C58D914'),
                    },
                };
            },
        },
    }), [theme]);

    const permissionsOptions = getFileAccess().map((entry: FileAccess) => ({
        value: entry.toString(),
        label:
      i18n[`types.permissions.${entry.toString().toLowerCase() as 'edit' | 'read'}`] ||
      entry.toString(),
    }));
    const [selectedUsers, setSelectedUsers] = useState<MattermostUser[]>([]);
    const [accessHeader, setAccessHeader] = useState<string>(i18n['permissions.loading']);

    useEffect(() => {
        const isChannel = window.location.href.split('/').includes('channels');
        if (!loading) {
            setAccessHeader(
                isChannel ? i18n['permissions.access_header_default'] : i18n['permissions.access_header'],
            );
        }
        return () => setSelectedUsers([]);
    }, [channel, loading, i18n]);

    const handleAddUsers = (): void => {
        if (selectedUsers.length > 0) {
            const contentSection = document.getElementById('scroller-dummy');
            setTimeout(() => contentSection?.scrollIntoView({behavior: 'smooth'}), 300);
            onAppendUsers(selectedUsers);
            setSelectedUsers([]);
        }
    };

    return (
        <div
            className={`filter-row onlyoffice-permissions__filter-row${!channel ? ' onlyoffice-permissions__filter-row--compact' : ''}`}
            data-theme={theme}
        >
            {channel && (
                <div className='col-xs-12 onlyoffice-permissions__user-column'>
                    <div className='onlyoffice-permissions__user-select-container'>
                        <div className='onlyoffice-permissions__async-select'>
                            <AsyncSelect
                                id='onlyoffice-permissions-select'
                                placeholder={i18n['permissions.modal_search_placeholder']}
                                loadingMessage={() => i18n['permissions.modal_search_loading_placeholder']}
                                noOptionsMessage={() => i18n['permissions.modal_search_no_options_placeholder']}
                                className='react-select-container'
                                classNamePrefix='react-select'
                                closeMenuOnSelect={false}
                                isMulti={true}
                                loadOptions={debounceUsersLoad(channel, fileInfo, users)}
                                onChange={(selected) => setSelectedUsers(selected as MattermostUser[])}
                                value={selectedUsers}
                                isDisabled={loading || !channel}
                                components={{
                                    DropdownIndicator: () => null,
                                    IndicatorSeparator: () => null,
                                }}
                                styles={styles.selectStyles}
                            />
                        </div>
                        <Button
                            className='btn btn-md btn-primary'
                            disabled={selectedUsers.length === 0 || loading}
                            onClick={handleAddUsers}
                        >
                            {i18n['permissions.modal_button_add']}
                        </Button>
                    </div>
                </div>
            )}
            <div className='col-sm-12 onlyoffice-permissions__bottom-section'>
                <span className='member-count pull-left onlyoffice-permissions__access-header'>
                    <span>{accessHeader}</span>
                </span>
                <div className='onlyoffice-permissions__permission-select'>
                    <Select
                        isSearchable={false}
                        value={{
                            value: wildcardAccess,
                            label:
                i18n[`types.permissions.${wildcardAccess.toLowerCase() as 'edit' | 'read'}`] ||
                wildcardAccess,
                        }}
                        options={permissionsOptions}
                        onChange={(selected) => onSetWildcardAccess(selected?.value)}
                        isDisabled={loading}
                        components={{
                            IndicatorSeparator: () => null,
                        }}
                        styles={styles.permissionSelectStyles}
                    />
                </div>
            </div>
        </div>
    );
};
