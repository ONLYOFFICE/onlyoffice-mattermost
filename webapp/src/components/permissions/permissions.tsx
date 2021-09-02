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

/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react/jsx-no-literals */
import React, {useState, useEffect} from 'react';
import {Dispatch} from 'redux';
import {FileInfo} from 'mattermost-redux/types/files';
import {Client4} from 'mattermost-redux/client';

import AsyncSelect from 'react-select/async';
import Select, {OptionTypeBase, OptionsType} from 'react-select';
import makeAnimated from 'react-select/animated';
import {Modal, Button} from 'react-bootstrap';

import {Channel} from 'mattermost-redux/types/channels';

import {apiGET, apiPOST} from 'api';
import {ONLYOFFICE_PLUGIN_API, ONLYOFFICE_PLUGIN_API_FILE_PERMISSIONS,
    ONLYOFFICE_PLUGIN_API_SET_FILE_PERMISSIONS, ONLYOFFICE_WILDCARD_USER} from 'utils';
import {debounce} from 'utils/lodash';
import {getTranslations} from 'utils/i18n';
import {FilePermissions, getFileAccess, getPermissionsMap, SubmitPermissionsPayload} from 'utils/file';
import {AutocompleteUser, mapUserToAutocompleteUser, sortAutocompleteUsers, User, getUniqueAutocompleteUsers} from 'utils/user';

import {UserRow} from './permissions_user_row';
import {PermissionsFooter} from './permissions_footer';
import {UserList} from './permissions_user_list';
import {PermissionsHeaderFilter} from './permissions_header_filter';

import 'public/scss/permissions.scss';

type PermissionsProps = {
    visible: boolean,
    close: () => (dispatch: Dispatch) => void,
    fileInfo: FileInfo,
};

const animatedComponents = makeAnimated();

//TODO: Refactoring
const Permissions: React.FC<PermissionsProps> = ({visible, close, fileInfo}: PermissionsProps) => {
    const i18n = getTranslations();
    const [isLoading, setIsLoading] = useState(true);
    const [allAccess, setAllAccess] = useState(FilePermissions.READ_ONLY.toString());
    const [current, setCurrent] = useState<AutocompleteUser[]>([]);
    const [users, setUsers] = useState<AutocompleteUser[]>([]);
    const [channel, setChannel] = useState<Channel>();
    const [accessHeaderText, setAccessHeaderText] = useState<string>(i18n['permissions.loading']);
    const permissionsMap = getPermissionsMap().map((entry: FilePermissions) => {
        return {
            value: entry.toString(),
            label: entry.toString(),
        };
    });

    useEffect(() => {
        if (visible) {
            (async () => {
                const response = await apiGET(ONLYOFFICE_PLUGIN_API + ONLYOFFICE_PLUGIN_API_FILE_PERMISSIONS + fileInfo.id);
                if (response[1].get('Channel-Type') === 'D') {
                    const arr = window.location.href.split('/');
                    setAccessHeaderText(`${i18n['permissions.access_header']} ${arr[arr.length - 1]}`);
                } else {
                    setAccessHeaderText(i18n['permissions.access_header_default']);

                    //TODO: Think of a better way to get current channel (fileInfo.post_id is not consistent)
                    const arr = window.location.href.split('/');
                    const team = await Client4.getTeamByName(arr[arr.length - 3]);
                    const chnl = await Client4.getChannelByName(team.id, arr[arr.length - 1]);
                    setChannel(chnl);
                }
                const resUsers: User[] = response[0];
                if (!resUsers) {
                    return;
                }
                const permissions: AutocompleteUser[] = [];
                // eslint-disable-next-line max-nested-callbacks
                resUsers.forEach((user: User) => {
                    const mappedUser = mapUserToAutocompleteUser(user);
                    if (user.id === ONLYOFFICE_WILDCARD_USER) {
                        setAllAccess(mappedUser.permissions);
                    } else {
                        permissions.push(mappedUser);
                    }
                });
                sortAutocompleteUsers(permissions);
                setUsers(permissions);
                setIsLoading(false);
            })();
        }
    }, [visible, fileInfo]);

    if (!visible) {
        return null;
    }

    const load = debounce((input: any, callback: any) => {
        if (!input) {
            return;
        }

        if (channel) {
            (async () => {
                let res = await Client4.searchUsers(input, {
                    in_channel_id: channel.id,
                    team_id: channel.team_id,
                });
                // eslint-disable-next-line max-nested-callbacks
                res = res.filter((user) => user.id !== fileInfo.user_id);
                const permissions = getUniqueAutocompleteUsers(res, users);
                callback(permissions);
            })();
        }
    }, 2000);

    const onChange = (value: OptionTypeBase | OptionsType<OptionTypeBase> | null) => {
        setCurrent((value as AutocompleteUser[]));
    };

    const onAllChange = (value: any) => {
        if (value.label) {
            setAllAccess(value.label);
        }
    };

    const onExit = () => {
        const modal = document.getElementById('onlyoffice-permissions-modal');
        const backdrop = modal?.previousElementSibling;

        // eslint-disable-next-line no-unused-expressions
        modal?.classList.remove('in');
        // eslint-disable-next-line no-unused-expressions
        backdrop?.classList.remove('in');

        setTimeout(() => {
            close();
            setAllAccess(FilePermissions.READ_ONLY.toString());
            setCurrent([]);
            setUsers([]);
            setAccessHeaderText(i18n['permissions.loading']);
            // eslint-disable-next-line no-undefined
            setChannel(undefined);
            setIsLoading(true);
        }, 300);
    };

    const onRemoveUser = (username: string) => {
        setUsers((prevUsers: AutocompleteUser[]) => prevUsers.filter((user: AutocompleteUser) => user.label !== username));
    };

    const onChangeUserPermissions = (username: string, newPermissions: string) => {
        setUsers((prevUsers: AutocompleteUser[]) => prevUsers.map((user: AutocompleteUser) => {
            if (user.label === username) {
                user.permissions = newPermissions;
            }
            return user;
        }));
    };

    const onSubmitPermissions = () => {
        const payload: SubmitPermissionsPayload[] = [];
        const allUsers: SubmitPermissionsPayload = {
            FileId: fileInfo.id,
            Username: '*',
            Id: '*',
            Permissions: FilePermissions.READ_ONLY.toString() === allAccess ? getFileAccess(FilePermissions.READ_ONLY) : getFileAccess(FilePermissions.EDIT_ONLY),
        };
        payload.push(allUsers);
        users.forEach((user: AutocompleteUser) => {
            payload.push({
                FileId: fileInfo.id,
                Username: user.label,
                Id: user.value,
                Permissions: FilePermissions.READ_ONLY.toString() === user.permissions ? getFileAccess(FilePermissions.READ_ONLY) : getFileAccess(FilePermissions.EDIT_ONLY),
            });
        });

        apiPOST(ONLYOFFICE_PLUGIN_API + ONLYOFFICE_PLUGIN_API_SET_FILE_PERMISSIONS, JSON.stringify(payload)).then(() => {
            onExit();
        }).catch();
    };

    return (
        <Modal
            show={visible}
            onHide={onExit}
            onExited={onExit}
            role='dialog'
            id='onlyoffice-permissions-modal'
        >
            <Modal.Header
                closeButton={true}
            >
                {`${i18n['permissions.modal_header']} ${fileInfo.name.split('.')[0]}`}
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={onExit}
                    disabled={isLoading}
                >
                    <span aria-hidden='true'>×</span>
                    <span className='sr-only'>Close</span>
                </button>
            </Modal.Header>
            <div
                className='onlyoffice-permissions-modal__body'
                style={channel ? {} : {maxHeight: '20rem'}}
            >
                <div
                    className='filtered-user-list'
                >
                    <div
                        className='filter-row'
                        style={channel ? {marginBottom: '1rem', marginTop: '1rem'} : {maxHeight: '10rem'}}
                    >
                        {channel && (
                            <PermissionsHeaderFilter>
                                <div style={{flexGrow: 1, marginRight: '2rem'}}>
                                    <AsyncSelect
                                        id='onlyoffice-permissions-select'
                                        className='react-select-container'
                                        classNamePrefix='react-select'
                                        closeMenuOnSelect={false}
                                        components={animatedComponents}
                                        isMulti={true}
                                        loadOptions={load}
                                        onChange={onChange}
                                        value={current}
                                        isDisabled={isLoading}
                                    />
                                </div>
                                <Button
                                    className='btn btn-md btn-primary'
                                    disabled={current.length === 0 || isLoading}
                                    onClick={() => {
                                        if (current) {
                                            const contentSection = document.getElementById('scroller-dummy');
                                            setTimeout(() => contentSection?.scrollIntoView({behavior: 'smooth'}), 300);
                                            setUsers((prevUsers: AutocompleteUser[]) => [...prevUsers, ...current]);
                                            setCurrent([]);
                                        }
                                    }}
                                >
                                    {i18n['permissions.modal_button_add']}
                                </Button>
                            </PermissionsHeaderFilter>
                        )}
                        <div
                            className='col-sm-12'
                            style={{marginTop: '2rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}
                        >
                            <span
                                style={{flexGrow: 2}}
                                className='member-count pull-left'
                            >
                                <span>{accessHeaderText}</span>
                            </span>
                            <div style={{marginRight: '2.5rem', marginLeft: '10px', width: '15rem'}}>
                                <Select
                                    isSearchable={false}
                                    value={{
                                        value: allAccess,
                                        label: allAccess,
                                    }}
                                    options={permissionsMap}
                                    onChange={onAllChange}
                                    isDisabled={isLoading}
                                />
                            </div>
                        </div>
                    </div>
                    {channel && (
                        <UserList>
                            {users.map((user: AutocompleteUser) => {
                                return (
                                    <UserRow
                                        key={user.value}
                                        user={user}
                                        changePermissions={onChangeUserPermissions}
                                        removeUser={onRemoveUser}
                                    />
                                );
                            })}
                            <div id='scroller-dummy'/>
                        </UserList>
                    )}
                    <PermissionsFooter>
                        <Button
                            className='btn btn-md'
                            style={{marginRight: '1rem', border: 'none'}}
                            onClick={onExit}
                            disabled={isLoading}
                        >
                            <span style={{color: '#2389D7'}}>{i18n['permissions.modal_button_cancel']}</span>
                        </Button>
                        <Button
                            className='btn btn-md btn-primary'
                            onClick={onSubmitPermissions}
                            disabled={isLoading}
                        >
                            {i18n['permissions.modal_button_save']}
                        </Button>
                    </PermissionsFooter>
                </div>
            </div>
        </Modal>
    );
};

export default Permissions;
