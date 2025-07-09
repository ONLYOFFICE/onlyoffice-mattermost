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

import {get, ONLYOFFICE_PLUGIN_CREATE, ONLYOFFICE_PLUGIN_GET_CODE, post} from 'api';
import React, {useState, useEffect} from 'react';
import {Modal} from 'react-bootstrap';
import {useSelector} from 'react-redux';
import type {Dispatch} from 'redux';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import ManagerActions from 'components/manager/ManagerActions';
import ManagerForm from 'components/manager/ManagerForm';
import ManagerHeader from 'components/manager/ManagerHeader';

import 'public/scss/manager.scss';

type Props = {
    visible: boolean;
    theme: string;
    darkTheme: string | undefined;
    close: () => (dispatch: Dispatch) => void;
};

const removeInAnimation = (): void => {
    const modal = document.getElementById('onlyoffice-manager-modal');
    const backdrop = modal?.previousElementSibling;
    modal?.classList.remove('in');
    backdrop?.classList.remove('in');
};

const getDefaultFileName = (type: string, i18n: Record<string, string>): string => {
    switch (type) {
        case 'pptx':
            return i18n['manager.default_name.presentation'];
        case 'xlsx':
            return i18n['manager.default_name.spreadsheet'];
        default:
            return i18n['manager.default_name.document'];
    }
};

export default function Manager({visible, theme, darkTheme, close}: Props) {
    const i18n = getTranslations();
    const channelId = useSelector(getCurrentChannelId);
    const [fileType, setFileType] = useState<string>('docx');
    const [fileName, setFileName] = useState<string>(getDefaultFileName('docx', i18n));
    const [error, setError] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);

    useEffect(() => {
        const currentDefaultName = getDefaultFileName(fileType, i18n);
        const otherDefaultNames = [
            i18n['manager.default_name.document'],
            i18n['manager.default_name.presentation'],
            i18n['manager.default_name.spreadsheet'],
        ];

        if (otherDefaultNames.includes(fileName))
            setFileName(currentDefaultName);
    }, [fileType, i18n]);

    if (!visible) {
        return null;
    }

    const handleExit = (): void => {
        removeInAnimation();
        setTimeout(() => close(), 300);
    };

    const handleCreate = async (): Promise<void> => {
        if (!fileName.trim()) {
            setError(i18n['manager.error_empty_name']);
            return;
        }

        if (fileName.length > 250) {
            setError(i18n['manager.error_name_too_long']);
            return;
        }

        setLoading(true);
        setError('');

        try {
            const code = await get<string>(ONLYOFFICE_PLUGIN_GET_CODE);
            await post(`${ONLYOFFICE_PLUGIN_CREATE}?code=${code}`, {
                channel_id: channelId,
                file_name: fileName,
                file_type: fileType,
            }, {
                credentials: 'include',
            });

            setFileName(getDefaultFileName('docx', i18n));
            setFileType('docx');
            handleExit();
        } catch (error) {
            setError(i18n['manager.error_create_failed']);
        } finally {
            setLoading(false);
        }
    };

    const handleFileNameChange = (value: string): void => {
        setFileName(value);
        if (!value.trim()) {
            setError(i18n['manager.error_empty_name']);
        } else if (value.length > 250) {
            setError(i18n['manager.error_name_too_long']);
        } else {
            setError('');
        }
    };

    return (
        <Modal
            show={visible}
            onHide={handleExit}
            onExited={handleExit}
            role='dialog'
            id='onlyoffice-manager-modal'
            data-theme={theme}
            data-dark-theme={darkTheme}
        >
            <ManagerHeader
                theme={theme}
                loading={loading}
                onClose={handleExit}
            />

            <div className='onlyoffice-manager-modal__body'>
                <div className='onlyoffice-manager__container'>
                    <ManagerForm
                        fileType={fileType}
                        fileName={fileName}
                        loading={loading}
                        error={error}
                        theme={theme}
                        darkTheme={darkTheme}
                        onFileTypeChange={setFileType}
                        onFileNameChange={handleFileNameChange}
                    />
                </div>

                <ManagerActions
                    loading={loading}
                    error={error}
                    onClose={handleExit}
                    onCreate={handleCreate}
                />
            </div>
        </Modal>
    );
}
