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

import ConverterActions from 'components/converter/ConverterActions';
import ConverterError from 'components/converter/ConverterError';
import ConverterFormatSelection from 'components/converter/ConverterFormatSelection';
import ConverterHeader from 'components/converter/ConverterHeader';
import ConverterInfo from 'components/converter/ConverterInfo';
import ConverterPasswordInput from 'components/converter/ConverterPasswordInput';

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
    const [selectedFormat, setSelectedFormat] = useState<'docx' | 'xlsx' | null>(null);

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
            } else if (response.error < 0) {
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
            <ConverterHeader
                theme={theme}
                onClose={handleClose}
            />
            <div className='onlyoffice-converter-modal__body'>
                <div className='onlyoffice-converter__container'>
                    <ConverterInfo/>
                    {needsPassword && (
                        <ConverterPasswordInput
                            password={password}
                            onPasswordChange={setPassword}
                        />
                    )}
                    {needsFormatSelection && (
                        <ConverterFormatSelection
                            selectedFormat={selectedFormat}
                            onFormatSelect={setSelectedFormat}
                        />
                    )}
                    <ConverterError error={error}/>
                </div>
                <ConverterActions
                    loading={loading}
                    needsPassword={needsPassword}
                    password={password}
                    needsFormatSelection={needsFormatSelection}
                    selectedFormat={selectedFormat}
                    onClose={handleClose}
                    onConvert={handleConvert}
                />
            </div>
        </Modal>
    );
}
