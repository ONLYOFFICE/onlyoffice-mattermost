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
    loading: boolean;
    onClose: () => void;
};

export default function ManagerHeader({theme, loading, onClose}: Props) {
    const i18n = getTranslations();

    return (
        <Modal.Header
            className='onlyoffice-manager-modal__header'
            data-theme={theme}
        >
            <span className='onlyoffice-manager-modal__header__text'>
                {i18n['manager.modal_header']}
            </span>
            <button
                type='button'
                className='close onlyoffice-manager-modal__header__close'
                aria-label='Close'
                onClick={onClose}
                disabled={loading}
            >
                <span aria-hidden='true'>{'×'}</span>
                <span className='sr-only'>{i18n['manager.cancel_button']}</span>
            </button>
        </Modal.Header>
    );
}
