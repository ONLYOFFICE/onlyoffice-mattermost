/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
import React, {useState, useEffect} from 'react';
import {Dispatch} from 'redux';
import {Modal} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';
import {Channel} from 'mattermost-redux/types/channels';
import {Client4} from 'mattermost-redux/client';

import {get, ONLYOFFICE_PLUGIN_PERMISSIONS} from 'api';

import {mapUsersToMattermostUsers, MattermostUser, OnlyofficeUser, sortMattermostUsers} from 'util/user';
import {FileAccess, getPermissionsTypeByPermissions} from 'util/permission';
import {ONLYOFFICE_WILDCARD_USER} from 'util/const';
import {getTranslations} from 'util/lang';
import {pipe} from 'util/func';

import {PermissionsFooter} from './PermissionsFooter';
import {PermissionsHeader} from './PermissionsHeader';
import {PermissionsList} from './PermissionsList';

import 'public/scss/permissions.scss';

type Props = {
    visible: boolean,
    close: () => (dispatch: Dispatch) => void,
    fileInfo: FileInfo
}

const removeInAnimation = () => {
    const modal = document.getElementById('onlyoffice-permissions-modal');
    const backdrop = modal?.previousElementSibling;
    // eslint-disable-next-line no-unused-expressions
    modal?.classList.remove('in');
    // eslint-disable-next-line no-unused-expressions
    backdrop?.classList.remove('in');
};

export default function OnlyofficeFilePermissions({visible, close, fileInfo}: Props) {
    const i18n = getTranslations();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);
    const [channel, setChannel] = useState<Channel | null>();
    const [users, setUsers] = useState<MattermostUser[]>([]);
    const [wildcardAccess, setWildcardAccess] = useState<string>(FileAccess.READ_ONLY);

    const onLoading = async () => {
        setChannel(null);
        const arr = window.location.href.split('/');
        try {
            if (arr.includes('channels')) {
                const team = await Client4.getTeamByName(arr[arr.length - 3]);
                const chnl = await Client4.getChannelByName(team.id, arr[arr.length - 1]);
                setChannel(chnl);
            }
            const response = await get<OnlyofficeUser[]>(`${ONLYOFFICE_PLUGIN_PERMISSIONS}?file=${fileInfo.id}`) || [];
            pipe<any>(getPermissionsTypeByPermissions, setWildcardAccess)(response.find((user) => user.id === ONLYOFFICE_WILDCARD_USER)?.permissions);
            pipe<any>(mapUsersToMattermostUsers, sortMattermostUsers, setUsers)(response);
        } catch (err) {
            setError(true);
        } finally {
            setLoading(false);
        }
    };

    const onExit = () => {
        removeInAnimation();
        setTimeout(() => {
            close();
        }, 300);
    };

    const onAppendUsers = (newUsers: MattermostUser[]) => {
        setUsers([...new Set([...users, ...newUsers])]);
    };

    const onRemoveUser = (username: string) => {
        const newUsers = users.filter((user) => user.label !== username);
        setUsers([...newUsers]);
    };

    const onChangeUserPermissions = (username: string, newPermissions: string) => {
        setUsers((prevUsers: MattermostUser[]) => prevUsers.map((user: MattermostUser) => {
            if (user.label === username) {
                user.fileAccess = newPermissions;
            }
            return user;
        }));
    };

    useEffect(() => {
        if (!visible) {
            return;
        }
        onLoading();
    }, [visible]);

    if (!visible) {
        return null;
    }

    return (
        <Modal
            show={visible}
            onHide={onExit}
            onExited={onExit}
            role='dialog'
            id='onlyoffice-permissions-modal'
        >
            <Modal.Header closeButton={true}>
                {`${i18n['permissions.modal_header']} ${fileInfo.name}`}
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={onExit}
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
                        onAppendUsers={onAppendUsers}
                        onSetWildcardAccess={setWildcardAccess}
                    />
                    {channel && (
                        <PermissionsList
                            users={users}
                            error={error}
                            onRemoveUser={onRemoveUser}
                            onChangeUserPermissions={onChangeUserPermissions}
                        />
                    )}
                    <PermissionsFooter
                        users={users}
                        onClose={onExit}
                        fileInfo={fileInfo}
                        loading={loading || error}
                        wildcardAccess={wildcardAccess}
                    />
                </div>
            </div>
        </Modal>
    );
}
