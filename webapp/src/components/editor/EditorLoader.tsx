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

import {ONLYOFFICE_CLOSE_EVENT, ONLYOFFICE_ERROR_EVENT, ONLYOFFICE_READY_EVENT} from 'util/const';
import {getTranslations} from 'util/lang';

import errorIcon from 'public/images/error.svg';
import React, {useEffect, useState} from 'react';

type Props = {
    theme: string;
};

export default function EditorLoader({theme}: Props) {
    const [error, setError] = useState(false);
    const [isVisible, setIsVisible] = useState(true);
    const i18n = getTranslations();

    const disableLoading = () => {
        setIsVisible(false);
    };

    const requestClose = () => {
        window.dispatchEvent(new Event(ONLYOFFICE_CLOSE_EVENT));
    };

    const trackError = () => setError(true);

    useEffect(() => {
        window.addEventListener(ONLYOFFICE_READY_EVENT, disableLoading);
        window.addEventListener(ONLYOFFICE_ERROR_EVENT, trackError);
        return () => {
            window.removeEventListener(ONLYOFFICE_READY_EVENT, disableLoading);
            window.removeEventListener(ONLYOFFICE_ERROR_EVENT, trackError);
        };
    }, []);

    if (!isVisible) {
        return null;
    }

    return (
        <div
            className='onlyoffice-editor__loader-container'
            data-theme={theme}
        >
            {!error && <div className='onlyoffice-editor__loader-icon'><div/><div/><div/></div>}
            {error && (
                <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center'}}>
                    <img
                        style={{width: '41px', height: '41px', marginBottom: '2rem'}}
                        src={errorIcon}
                    />
                    <span className='onlyoffice-editor__loader_error'>{i18n['editor.events.error']}</span>
                </div>
            )}
            <button
                className='onlyoffice-editor__loader-btn'
                onClick={requestClose}
            >
                {i18n['editor.close_button']}
            </button>
        </div>
    );
}
