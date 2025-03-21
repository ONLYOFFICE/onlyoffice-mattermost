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

import {ONLYOFFICE_WILDCARD_USER} from 'util/const';
import {pipe} from 'util/func';
import {getTranslations} from 'util/lang';
import {FileAccess, getPermissionsTypeByPermissions} from 'util/permission';
import type {MattermostUser, OnlyofficeUser} from 'util/user';
import {mapUsersToMattermostUsers, sortMattermostUsers} from 'util/user';

import {get, ONLYOFFICE_PLUGIN_PERMISSIONS} from 'api';
import React, {useState, useEffect} from 'react';
import {Modal} from 'react-bootstrap';
import type {Dispatch} from 'redux';

import {Client4} from 'mattermost-redux/client';
import type {Channel} from 'mattermost-redux/types/channels';
import type {FileInfo} from 'mattermost-redux/types/files';

import {PermissionsFooter} from './PermissionsFooter';
import {PermissionsHeader} from './PermissionsHeader';
import {PermissionsList} from './PermissionsList';

import 'public/scss/permissions.scss';

type Props = {
    visible: boolean;
    close: () => (dispatch: Dispatch) => void;
    fileInfo: FileInfo;
};

const removeInAnimation = (): void => {
    const modal = document.getElementById('onlyoffice-permissions-modal');
    const backdrop = modal?.previousElementSibling;
    modal?.classList.remove('in');
    backdrop?.classList.remove('in');
};

export default function OnlyofficeFilePermissions({visible, close, fileInfo}: Props) {
    const i18n = getTranslations();
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<boolean>(false);
    const [channel, setChannel] = useState<Channel | null>(null);
    const [users, setUsers] = useState<MattermostUser[]>([]);
    const [wildcardAccess, setWildcardAccess] = useState<string>(FileAccess.READ_ONLY);

    const fetchData = async (): Promise<void> => {
        setChannel(null);
        const urlParts = window.location.href.split('/');
        try {
            if (urlParts.includes('channels')) {
                const teamName = urlParts[urlParts.length - 3];
                const channelName = urlParts[urlParts.length - 1];
                const team = await Client4.getTeamByName(teamName);
                const chnl = await Client4.getChannelByName(team.id, channelName);
                setChannel(chnl);
            }
            const response = (await get<OnlyofficeUser[]>(
                `${ONLYOFFICE_PLUGIN_PERMISSIONS}?file=${fileInfo.id}`,
            )) || [];
            pipe<any>(getPermissionsTypeByPermissions, setWildcardAccess)(
                response.find((user) => user.id === ONLYOFFICE_WILDCARD_USER)?.permissions,
            );
            pipe<any>(mapUsersToMattermostUsers, sortMattermostUsers, setUsers)(response);
        } catch (err) {
            setError(true);
        } finally {
            setLoading(false);
        }
    };

    const handleExit = (): void => {
        removeInAnimation();
        setTimeout(() => close(), 300);
    };

    const handleAppendUsers = (newUsers: MattermostUser[]): void => {
        setUsers((prevUsers) => {
            const allUsers = [...prevUsers, ...newUsers];
            return [...new Set(allUsers)];
        });
    };

    const handleRemoveUser = (username: string): void => {
        setUsers((prevUsers) => prevUsers.filter((user) => user.label !== username));
    };

    const handleChangeUserPermissions = (username: string, newPermissions: string): void => {
        setUsers((prevUsers) =>
            prevUsers.map((user) =>
                (user.label === username ? {...user, fileAccess: newPermissions} : user),
            ),
        );
    };

    useEffect(() => {
        if (visible) {
            fetchData();
        }
    }, [visible]);

    if (!visible) {
        return null;
    }

    return (
        <Modal
            show={visible}
            onHide={handleExit}
            onExited={handleExit}
            role='dialog'
            id='onlyoffice-permissions-modal'
        >
            <Modal.Header className='onlyoffice-permissions-modal__header'>
                <span style={{fontWeight: 600}}>{`${i18n['permissions.modal_header']} ${fileInfo.name}`}</span>
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={handleExit}
                    disabled={loading}
                >
                    <span aria-hidden='true'>{'×'}</span>
                    <span className='sr-only'>{'Close'}</span>
                </button>
            </Modal.Header>
            <div
                className='onlyoffice-permissions-modal__body'
                style={channel ? {} : {maxHeight: '20rem'}}
            >
                <div className='filtered-user-list'>
                    <PermissionsHeader
                        fileInfo={fileInfo}
                        channel={channel}
                        loading={loading}
                        wildcardAccess={wildcardAccess}
                        users={users}
                        onAppendUsers={handleAppendUsers}
                        onSetWildcardAccess={setWildcardAccess}
                    />
                    {channel && (
                        <PermissionsList
                            users={users}
                            error={error}
                            onRemoveUser={handleRemoveUser}
                            onChangeUserPermissions={handleChangeUserPermissions}
                        />
                    )}
                    <PermissionsFooter
                        users={users}
                        onClose={handleExit}
                        fileInfo={fileInfo}
                        loading={loading || error}
                        wildcardAccess={wildcardAccess}
                    />
                </div>
            </div>
        </Modal>
    );
}

