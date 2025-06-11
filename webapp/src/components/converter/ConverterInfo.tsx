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

export default function ConverterInfo() {
    const i18n = getTranslations();

    return (
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
    );
}
