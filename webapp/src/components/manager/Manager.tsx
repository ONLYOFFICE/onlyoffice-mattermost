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

import React, {useState} from 'react';
import {Modal} from 'react-bootstrap';
import {useSelector} from 'react-redux';
import type {Dispatch} from 'redux';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {get, ONLYOFFICE_PLUGIN_CREATE, ONLYOFFICE_PLUGIN_GET_CODE, post} from 'api';

import {getTranslations} from 'util/lang';

import ManagerHeader from 'components/manager/ManagerHeader';
import ManagerForm from 'components/manager/ManagerForm';
import ManagerError from 'components/manager/ManagerError';
import ManagerActions from 'components/manager/ManagerActions';

import 'public/scss/manager.scss';

type Props = {
    visible: boolean;
    theme: string;
    close: () => (dispatch: Dispatch) => void;
};

export default function Manager({visible, theme, close}: Props) {
    const i18n = getTranslations();
    const channelId = useSelector(getCurrentChannelId);
    const [fileType, setFileType] = useState<string>('docx');
    const [fileName, setFileName] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);

    if (!visible) {
        return null;
    }

    const handleCreate = async (): Promise<void> => {
        setError('');

        if (!fileName.trim()) {
            setError(i18n['manager.error_empty_name']);
            return;
        }

        setLoading(true);

        try {
            const code = await get<string>(ONLYOFFICE_PLUGIN_GET_CODE);
            await post(`${ONLYOFFICE_PLUGIN_CREATE}?code=${code}`, {
                channel_id: channelId,
                file_name: fileName,
                file_type: fileType,
            }, {
                credentials: 'include',
            });

            setFileName('');
            setFileType('docx');
            close();
        } catch (error) {
            setError(i18n['manager.error_create_failed']);
        } finally {
            setLoading(false);
        }
    };

    const handleClose = (): void => {
        if (!loading) {
            close();
        }
    };

    return (
        <Modal
            show={visible}
            onHide={handleClose}
            onExited={handleClose}
            role='dialog'
            id='onlyoffice-manager-modal'
            data-theme={theme}
        >
            <ManagerHeader
                theme={theme}
                loading={loading}
                onClose={handleClose}
            />

            <div className='onlyoffice-manager-modal__body'>
                <div className='onlyoffice-manager__container'>
                    <ManagerForm
                        fileType={fileType}
                        fileName={fileName}
                        loading={loading}
                        onFileTypeChange={setFileType}
                        onFileNameChange={setFileName}
                    />

                    <ManagerError error={error} />
                </div>

                <ManagerActions
                    loading={loading}
                    onClose={handleClose}
                    onCreate={handleCreate}
                />
            </div>
        </Modal>
    );
}
