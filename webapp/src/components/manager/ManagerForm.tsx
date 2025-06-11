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

import React from 'react';

type Props = {
    fileType: string;
    fileName: string;
    loading: boolean;
    error: string;
    theme: string;
    darkTheme: string | undefined;
    onFileTypeChange: (fileType: string) => void;
    onFileNameChange: (fileName: string) => void;
};

type TranslationType = {
    'manager.file_type.document': string;
    'manager.file_type.spreadsheet': string;
    'manager.file_type.presentation': string;
    'manager.file_name_label': string;
    [key: string]: string;
};

const types = [
    {label: 'manager.file_type.document' as const, value: 'docx'},
    {label: 'manager.file_type.spreadsheet' as const, value: 'xlsx'},
    {label: 'manager.file_type.presentation' as const, value: 'pptx'},
];

export default function ManagerForm({
    fileType,
    fileName,
    loading,
    error,
    theme,
    darkTheme,
    onFileTypeChange,
    onFileNameChange,
}: Props) {
    const i18n = getTranslations() as TranslationType;

    return (
        <>
            <div className='onlyoffice-manager__form-row'>
                <input
                    type='text'
                    value={fileName}
                    onChange={(e) => onFileNameChange(e.target.value)}
                    disabled={loading}
                    className='onlyoffice-manager__input onlyoffice-manager__text-input'
                    placeholder={i18n['manager.file_name_label']}
                    data-theme={theme}
                    data-dark-theme={darkTheme}
                />
                {error && !fileName.trim() && <div className='onlyoffice-manager__error'>{error}</div>}
            </div>

            <div className='onlyoffice-manager__form-row'>
                <div className='onlyoffice-manager__select-container'>
                    <select
                        value={fileType}
                        onChange={(e) => onFileTypeChange(e.target.value)}
                        disabled={loading}
                        className='onlyoffice-manager__select'
                        data-theme={theme}
                        data-dark-theme={darkTheme}
                    >
                        {types.map((type) => (
                            <option
                                key={type.value}
                                value={type.value}
                            >
                                {i18n[type.label]}
                            </option>
                        ))}
                    </select>
                </div>
            </div>
        </>
    );
}
