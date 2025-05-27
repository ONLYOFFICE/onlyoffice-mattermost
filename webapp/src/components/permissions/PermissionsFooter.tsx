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
import {getTranslations} from 'util/lang';
import type {SubmitPermissionsRequest} from 'util/permission';
import {getFilePermissions} from 'util/permission';
import type {MattermostUser} from 'util/user';

import {get, post, ONLYOFFICE_PLUGIN_GET_CODE, ONLYOFFICE_PLUGIN_PERMISSIONS} from 'api';
import React from 'react';
import {Button} from 'react-bootstrap';

import type {FileInfo} from 'mattermost-redux/types/files';

type Props = {
    fileInfo: FileInfo;
    loading: boolean;
    users: MattermostUser[];
    wildcardAccess: string;
    onClose: () => void;
    theme: string;
};

const styles = {
    container: {
        display: 'flex',
        justifyContent: 'flex-end',
    },
    cancelButton: (theme: string) => ({
        marginRight: '1rem',
        border: 'none',
        backgroundColor: theme === 'dark' ? '#1b1d22' : undefined,
        color: theme === 'dark' ? '#ffffff' : undefined,
    }),
    cancelText: {
        color: 'var(--button-bg)',
    },
    saveButton: (theme: string) => ({
        backgroundColor: theme === 'dark' ? '#166de0' : undefined,
        color: theme === 'dark' ? '#ffffff' : undefined,
    }),
};

export const PermissionsFooter: React.FC<Props> = ({
    fileInfo,
    loading,
    users,
    wildcardAccess,
    onClose,
    theme,
}) => {
    const i18n = getTranslations();

    const handleSubmit = async (): Promise<void> => {
        const submitRequests: SubmitPermissionsRequest[] = [
            {
                fileID: fileInfo.id,
                userID: ONLYOFFICE_WILDCARD_USER,
                username: ONLYOFFICE_WILDCARD_USER,
                permissions: getFilePermissions(wildcardAccess),
            },
            ...users.map((user) => ({
                fileID: fileInfo.id,
                userID: user.value,
                username: user.label,
                permissions: getFilePermissions(user.fileAccess),
            })),
        ];

        try {
            // TODO: Handle too many permission entries if needed.
            if (submitRequests.length <= 25) {
                const code = await get<string>(ONLYOFFICE_PLUGIN_GET_CODE);
                await post<SubmitPermissionsRequest[], void>(
                    `${ONLYOFFICE_PLUGIN_PERMISSIONS}?code=${code}`,
                    submitRequests,
                );
            }
        } finally {
            onClose();
        }
    };

    return (
        <div
            className='filter-controls'
            style={styles.container}
            data-theme={theme}
        >
            <Button
                className='btn btn-md'
                style={styles.cancelButton(theme)}
                disabled={loading}
                onClick={onClose}
            >
                <span style={styles.cancelText}>
                    {i18n['permissions.modal_button_cancel']}
                </span>
            </Button>
            <Button
                className='btn btn-md btn-primary'
                onClick={handleSubmit}
                disabled={loading}
                style={styles.saveButton(theme)}
            >
                {i18n['permissions.modal_button_save']}
            </Button>
        </div>
    );
};

