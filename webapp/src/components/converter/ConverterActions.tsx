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
    loading: boolean;
    needsPassword: boolean;
    password: string;
    needsFormatSelection: boolean;
    selectedFormat: 'docx' | 'xlsx' | null;
    onClose: () => void;
    onConvert: () => void;
};

export default function ConverterActions({
    loading,
    needsPassword,
    password,
    needsFormatSelection,
    selectedFormat,
    onClose,
    onConvert,
}: Props) {
    const i18n = getTranslations();

    const isConvertDisabled = loading || 
        (needsPassword && !password) || 
        (needsFormatSelection && !selectedFormat);

    return (
        <div className='onlyoffice-converter__actions'>
            <button
                className='btn btn-secondary onlyoffice-converter__button onlyoffice-converter__cancel-button'
                onClick={onClose}
                disabled={loading}
            >
                {i18n['converter.cancel_button'] || 'Cancel'}
            </button>
            <button
                className='btn btn-primary onlyoffice-converter__button'
                onClick={onConvert}
                disabled={isConvertDisabled}
            >
                {loading ?
                    (i18n['converter.converting_button'] || 'Converting...') :
                    (i18n['converter.convert_button'] || 'Convert')
                }
            </button>
        </div>
    );
} 