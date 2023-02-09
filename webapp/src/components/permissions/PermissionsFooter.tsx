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
import React from 'react';
import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';

import {get, ONLYOFFICE_PLUGIN_GET_CODE, ONLYOFFICE_PLUGIN_PERMISSIONS, post} from 'api';

import {MattermostUser} from 'util/user';
import {ONLYOFFICE_WILDCARD_USER} from 'util/const';
import {getTranslations} from 'util/lang';
import {getFilePermissions, SubmitPermissionsRequest} from 'util/permission';

type Props = {
    fileInfo: FileInfo,
    loading: boolean,
    users: MattermostUser[],
    wildcardAccess: string,
    onClose: () => void,
};

const onSubmit = async (props: Props) => {
    const requestBody: SubmitPermissionsRequest[] = [];
    const wildcardBody: SubmitPermissionsRequest = {
        fileID: props.fileInfo.id,
        userID: ONLYOFFICE_WILDCARD_USER,
        username: ONLYOFFICE_WILDCARD_USER,
        permissions: getFilePermissions(props.wildcardAccess),
    };

    requestBody.push(wildcardBody);
    props.users.forEach((user) => {
        requestBody.push({
            fileID: props.fileInfo.id,
            userID: user.value,
            username: user.label,
            permissions: getFilePermissions(user.fileAccess),
        });
    });

    try {
        //TODO: Handle too many permission entries
        if (requestBody.length <= 25) {
            const code = await get<string>(ONLYOFFICE_PLUGIN_GET_CODE);
            await post<SubmitPermissionsRequest[], void>(`${ONLYOFFICE_PLUGIN_PERMISSIONS}?code=${code}`, requestBody);
        }
    } finally {
        props.onClose();
    }
};

export const PermissionsFooter = (props: Props) => {
    const i18n = getTranslations();
    return (
        <div
            className='filter-controls'
            style={{display: 'flex', justifyContent: ' flex-end', padding: 0, margin: '1rem', maxHeight: '4rem'}}
        >
            <Button
                className='btn btn-md'
                style={{marginRight: '1rem', border: 'none'}}
                disabled={props.loading}
                onClick={props.onClose}
            >
                <span style={{color: 'var(--button-bg)'}}>{i18n['permissions.modal_button_cancel']}</span>
            </Button>
            <Button
                className='btn btn-md btn-primary'
                onClick={() => onSubmit(props)}
                disabled={props.loading}
            >
                {i18n['permissions.modal_button_save']}
            </Button>
        </div>
    );
};
