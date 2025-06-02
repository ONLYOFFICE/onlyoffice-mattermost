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
    selectedFormat: 'docx' | 'xlsx' | null;
    onFormatSelect: (format: 'docx' | 'xlsx') => void;
};

export default function ConverterFormatSelection({selectedFormat, onFormatSelect}: Props) {
    const i18n = getTranslations();

    return (
        <div className='onlyoffice-converter__format-section'>
            <div className='onlyoffice-converter__format-title'>
                {i18n['converter.select_format'] || 'Select output format:'}
            </div>
            <div className='onlyoffice-converter__format-buttons'>
                <button
                    className={`onlyoffice-converter__format-button document ${selectedFormat === 'docx' ? 'selected' : ''}`}
                    onClick={() => onFormatSelect('docx')}
                >
                    {i18n['converter.format_document'] || 'Document'}
                </button>
                <button
                    className={`onlyoffice-converter__format-button cell ${selectedFormat === 'xlsx' ? 'selected' : ''}`}
                    onClick={() => onFormatSelect('xlsx')}
                >
                    {i18n['converter.format_cell'] || 'Spreadsheet'}
                </button>
            </div>
        </div>
    );
} 