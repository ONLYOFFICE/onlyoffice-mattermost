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
import {Modal} from 'react-bootstrap';

type Props = {
    theme: string;
    onClose: () => void;
};

export default function ConverterHeader({theme, onClose}: Props) {
    const i18n = getTranslations();

    return (
        <Modal.Header
            className='onlyoffice-converter-modal__header'
            data-theme={theme}
        >
            <span
                className='onlyoffice-converter-modal__header__text'
                style={{fontWeight: 600}}
            >
                {i18n['converter.modal_header'] || 'Convert File'}
            </span>
            <button
                type='button'
                className='close onlyoffice-converter-modal__header__close'
                aria-label='Close'
                onClick={onClose}
            >
                <span aria-hidden='true'>{'×'}</span>
                <span className='sr-only'>{i18n['converter.cancel_button'] || 'Cancel'}</span>
            </button>
        </Modal.Header>
    );
}
