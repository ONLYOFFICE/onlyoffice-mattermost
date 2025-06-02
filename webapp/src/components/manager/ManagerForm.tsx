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

import React from 'react';

import {getTranslations} from 'util/lang';

type Props = {
    fileType: string;
    fileName: string;
    loading: boolean;
    onFileTypeChange: (fileType: string) => void;
    onFileNameChange: (fileName: string) => void;
};

const types = [
    {label: 'DOCX', value: 'docx'},
    {label: 'XLSX', value: 'xlsx'},
    {label: 'PPTX', value: 'pptx'},
];

export default function ManagerForm({
    fileType,
    fileName,
    loading,
    onFileTypeChange,
    onFileNameChange,
}: Props) {
    const i18n = getTranslations();

    return (
        <>
            <div className='onlyoffice-manager__form-row'>
                <label className='onlyoffice-manager__label'>
                    {i18n['manager.file_type_label']}
                </label>
                <select
                    value={fileType}
                    onChange={(e) => onFileTypeChange(e.target.value)}
                    disabled={loading}
                    className='onlyoffice-manager__input onlyoffice-manager__select'
                >
                    {types.map((type) => (
                        <option
                            key={type.value}
                            value={type.value}
                        >
                            {type.label}
                        </option>
                    ))}
                </select>
            </div>

            <div className='onlyoffice-manager__form-row'>
                <label className='onlyoffice-manager__label'>
                    {i18n['manager.file_name_label']}
                </label>
                <input
                    type='text'
                    value={fileName}
                    onChange={(e) => onFileNameChange(e.target.value)}
                    disabled={loading}
                    className='onlyoffice-manager__input onlyoffice-manager__text-input'
                />
            </div>
        </>
    );
} 