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
    password: string;
    onPasswordChange: (password: string) => void;
};

export default function ConverterPasswordInput({password, onPasswordChange}: Props) {
    const i18n = getTranslations();

    return (
        <div className='onlyoffice-converter__password-section'>
            <div className='onlyoffice-converter__password-input-container'>
                <input
                    type='password'
                    className='onlyoffice-converter__password-input form-control'
                    placeholder={i18n['converter.password_placeholder'] || 'Enter file password'}
                    value={password}
                    onChange={(e) => onPasswordChange(e.target.value)}
                    required
                />
            </div>
        </div>
    );
} 