/**
 *
 * (c) Copyright Ascensio System SIA 2022
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
import React, {useEffect, useState} from 'react';

import {ONLYOFFICE_CLOSE_EVENT, ONLYOFFICE_ERROR_EVENT, ONLYOFFICE_READY_EVENT} from 'util/const';
import {getTranslations} from 'util/lang';

export default function EditorLoader() {
    const [error, setError] = useState(false);
    const i18n = getTranslations();

    const disableLoading = () => {
        const container = document.getElementsByClassName('onlyoffice-editor__loader-container').item(0);
        if (container) {
            container.classList.add('onlyoffice-editor__loader-container_hidden');
        }
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

    return (
        <div className='onlyoffice-editor__loader-container'>
            {!error && <div className='onlyoffice-editor__loader-icon'><div/><div/><div/></div>}
            {error && <span className='onlyoffice-editor__loader_error'>{i18n['editor.events.error']}</span>}
            <button
                className='onlyoffice-editor__loader-btn'
                onClick={requestClose}
            >
                {i18n['editor.close_button']}
            </button>
        </div>
    );
}
