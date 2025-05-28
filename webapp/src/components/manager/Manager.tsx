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
import type {Dispatch} from 'redux';

import React, {useState} from 'react';
import { Modal } from 'react-bootstrap';
import {get, ONLYOFFICE_PLUGIN_CREATE, ONLYOFFICE_PLUGIN_GET_CODE, post} from 'api';

import {useSelector} from 'react-redux';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import 'public/scss/manager.scss';

type Props = {
    visible: boolean;
    close: () => (dispatch: Dispatch) => void;
};

const types = [
    {label: 'DOCX', value: 'docx'},
    {label: 'XLSX', value: 'xlsx'},
    {label: 'PPTX', value: 'pptx'},
];

export default function Manager({visible, close}: Props) {
    const i18n = getTranslations() as any;
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
        if (!loading) close();
    };

    return (
        <Modal
            show={visible}
            onHide={handleClose}
            onExited={handleClose}
            role="dialog"
            id="onlyoffice-manager-modal"
            data-theme="dark"
        >
            <Modal.Header className="onlyoffice-manager-modal__header" data-theme="dark">
                <span style={{fontWeight: 600}}>
                    {i18n['manager.modal_header']}
                </span>
                <button
                    type="button"
                    className="close"
                    aria-label="Close"
                    onClick={handleClose}
                    disabled={loading}
                >
                    <span aria-hidden="true">{'×'}</span>
                    <span className="sr-only">{i18n['manager.cancel_button']}</span>
                </button>
            </Modal.Header>
            
            <div className="onlyoffice-manager-modal__body">
                <div className="onlyoffice-manager__container">
                    <div className="onlyoffice-manager__form-row">
                        <label className="onlyoffice-manager__label">
                            {i18n['manager.file_type_label']}
                        </label>
                        <select
                            value={fileType}
                            onChange={(e) => setFileType(e.target.value)}
                            disabled={loading}
                            className="onlyoffice-manager__input onlyoffice-manager__select"
                        >
                            {types.map((type) => (
                                <option key={type.value} value={type.value}>
                                    {type.label}
                                </option>
                            ))}
                        </select>
                    </div>
                    
                    <div className="onlyoffice-manager__form-row">
                        <label className="onlyoffice-manager__label">
                            {i18n['manager.file_name_label']}
                        </label>
                        <input
                            type="text"
                            value={fileName}
                            onChange={(e) => setFileName(e.target.value)}
                            disabled={loading}
                            placeholder={i18n['manager.file_name_placeholder'] || ''}
                            className="onlyoffice-manager__input onlyoffice-manager__text-input"
                        />
                    </div>
                    
                    {error && (
                        <div className="onlyoffice-manager__error">
                            {error}
                        </div>
                    )}
                </div>
                
                <div className="onlyoffice-manager__actions">
                    <button
                        className="btn btn-secondary onlyoffice-manager__button onlyoffice-manager__cancel-button"
                        onClick={handleClose}
                        disabled={loading}
                    >
                        {i18n['manager.cancel_button']}
                    </button>
                    <button
                        className="btn btn-primary onlyoffice-manager__button"
                        onClick={handleCreate}
                        disabled={loading}
                    >
                        {loading ? i18n['manager.creating_button'] || 'Creating...' : i18n['manager.create_button']}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
