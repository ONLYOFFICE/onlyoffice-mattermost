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

import {get, post, ONLYOFFICE_PLUGIN_GET_CODE, ONLYOFFICE_PLUGIN_CONVERT} from 'api';
import React, {useState} from 'react';
import {Modal} from 'react-bootstrap';
import type {Dispatch} from 'redux';

import type {FileInfo} from 'mattermost-redux/types/files';

import 'public/scss/converter.scss';

type Props = {
    visible: boolean;
    fileInfo: FileInfo;
    theme: string;
    close: () => (dispatch: Dispatch) => void;
};

export default function Converter({visible, fileInfo, theme, close}: Props) {
    const i18n = getTranslations() as any;
    const [error, setError] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);
    const [needsPassword, setNeedsPassword] = useState<boolean>(false);
    const [password, setPassword] = useState<string>('');
    const [needsFormatSelection, setNeedsFormatSelection] = useState<boolean>(false);
    const [selectedFormat, setSelectedFormat] = useState<'word' | 'cell' | null>(null);

    if (!visible) {
        return null;
    }

    const handleClose = (): void => {
        setNeedsPassword(false);
        setPassword('');
        setError('');
        setNeedsFormatSelection(false);
        setSelectedFormat(null);
        close();
    };

    const handleConvert = async (): Promise<void> => {
        setError('');
        setLoading(true);

        try {
            const code = await get<string>(ONLYOFFICE_PLUGIN_GET_CODE);
            const payload = {
                file_id: fileInfo.id,
                ...(needsPassword && {password}),
                ...(needsFormatSelection && {output_type: selectedFormat}),
            };
            const response = await post<{file_id: string}, {error: number}>(`${ONLYOFFICE_PLUGIN_CONVERT}?code=${code}`, payload, {
                credentials: 'include',
            });

            if (response.error === -5) {
                setNeedsPassword(true);
                setError(i18n['converter.error_password_required'] || 'This file is password protected. Please enter the password.');
            } else if (response.error === -9) {
                setNeedsFormatSelection(true);
                setError(i18n['converter.error_format_required'] || 'Please select the output format for conversion.');
            } else if (response.error !== 0) {
                throw new Error('Failed to convert file. Please try again.');
            } else {
                setNeedsPassword(false);
                setPassword('');
                setNeedsFormatSelection(false);
                setSelectedFormat(null);
                close();
            }
        } catch (error: any) {
            setError(i18n['converter.error_convert_failed'] || 'Failed to convert file. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            show={visible}
            onHide={handleClose}
            onExited={handleClose}
            role='dialog'
            id='onlyoffice-converter-modal'
            data-theme={theme}
        >
            <Modal.Header
                className='onlyoffice-converter-modal__header'
                data-theme={theme}
            >
                <span style={{fontWeight: 600}}>
                    {i18n['converter.modal_header'] || 'Convert File'}
                </span>
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={handleClose}
                >
                    <span aria-hidden='true'>{'×'}</span>
                    <span className='sr-only'>{i18n['converter.cancel_button'] || 'Cancel'}</span>
                </button>
            </Modal.Header>

            <div className='onlyoffice-converter-modal__body'>
                <div className='onlyoffice-converter__container'>
                    <div className='onlyoffice-converter__conversion-info'>
                        <div className='onlyoffice-converter__info-section'>
                            <h4>{i18n['converter.conversion_info_title'] || 'File Conversion'}</h4>
                            <div className='onlyoffice-converter__info-description'>
                                <p>{i18n['converter.conversion_description'] || 'This will convert your file to an OOXML format that is fully compatible with ONLYOFFICE editors.'}</p>
                                <ul className='onlyoffice-converter__conversion-list'>
                                    <li>{i18n['converter.word_conversion'] || 'Word documents (.doc, .odt, .rtf, etc.) → DOCX'}</li>
                                    <li>{i18n['converter.excel_conversion'] || 'Spreadsheets (.xls, .ods, .csv, etc.) → XLSX'}</li>
                                    <li>{i18n['converter.powerpoint_conversion'] || 'Presentations (.ppt, .odp, etc.) → PPTX'}</li>
                                </ul>
                                <div className='onlyoffice-converter__info-note'>
                                    <span className='onlyoffice-converter__note-text'>
                                        {i18n['converter.conversion_note'] || 'The converted file will be saved as a new attachment in this channel.'}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {needsPassword && (
                        <div className='onlyoffice-converter__password-section'>
                            <div className='onlyoffice-converter__password-input-container'>
                                <input
                                    type='password'
                                    className='onlyoffice-converter__password-input form-control'
                                    placeholder={i18n['converter.password_placeholder'] || 'Enter file password'}
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                                />
                            </div>
                        </div>
                    )}

                    {needsFormatSelection && (
                        <div className='onlyoffice-converter__format-section'>
                            <div className='onlyoffice-converter__format-title'>
                                {i18n['converter.select_format'] || 'Select output format:'}
                            </div>
                            <div className='onlyoffice-converter__format-buttons'>
                                <button
                                    className={`onlyoffice-converter__format-button document ${selectedFormat === 'word' ? 'selected' : ''}`}
                                    onClick={() => setSelectedFormat('word')}
                                >
                                    {i18n['converter.format_document'] || 'Document'}
                                </button>
                                <button
                                    className={`onlyoffice-converter__format-button cell ${selectedFormat === 'cell' ? 'selected' : ''}`}
                                    onClick={() => setSelectedFormat('cell')}
                                >
                                    {i18n['converter.format_cell'] || 'Spreadsheet'}
                                </button>
                            </div>
                        </div>
                    )}

                    {error && (
                        <div className='onlyoffice-converter__error'>
                            {error}
                        </div>
                    )}
                </div>

                <div className='onlyoffice-converter__actions'>
                    <button
                        className='btn btn-secondary onlyoffice-converter__button onlyoffice-converter__cancel-button'
                        onClick={handleClose}
                        disabled={loading}
                    >
                        {i18n['converter.cancel_button'] || 'Cancel'}
                    </button>
                    <button
                        className='btn btn-primary onlyoffice-converter__button'
                        onClick={handleConvert}
                        disabled={loading || (needsPassword && !password) || (needsFormatSelection && !selectedFormat)}
                    >
                        {loading ?
                            (i18n['converter.converting_button'] || 'Converting...') :
                            (i18n['converter.convert_button'] || 'Convert')
                        }
                    </button>
                </div>
            </div>
        </Modal>
    );
}
