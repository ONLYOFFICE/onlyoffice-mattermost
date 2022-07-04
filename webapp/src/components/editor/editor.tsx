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
import React, {useCallback, useEffect} from 'react';
import {Dispatch} from 'redux';
import {FileInfo} from 'mattermost-redux/types/files';

import {ONLYOFFICE_CLOSE_EVENT, ONLYOFFICE_PLUGIN_API} from 'util/const';

import EditorLoader from './EditorLoader';

type Props = {
    visible: boolean,
    fileInfo?: FileInfo,
    close: () => (dispatch: Dispatch) => void,
};

export default function Editor({visible, close, fileInfo}: Props) {
    const lang = localStorage.getItem('onlyoffice_locale') || 'en';
    const handleClose = useCallback(() => {
        if (!visible) {
            return;
        }
        const editorBackdrop = document.getElementById('editor-backdrop');

        if (editorBackdrop) {
            editorBackdrop.classList.add('onlyoffice-modal__backdrop_hide');
        }

        setTimeout(() => close(), 280);
    }, [close, visible]);

    useEffect(() => {
        window.addEventListener(ONLYOFFICE_CLOSE_EVENT, handleClose);
        return () => window.removeEventListener(ONLYOFFICE_CLOSE_EVENT, handleClose);
    }, [handleClose]);

    return (
        <>
            {visible && (
                <div
                    id='editor-backdrop'
                    className='onlyoffice-modal__backdrop'
                >
                    <EditorLoader/>
                    <iframe
                        src={`${ONLYOFFICE_PLUGIN_API}/editor?file=${fileInfo?.id}&lang=${lang}`}
                        className='onlyoffice-modal__frame'
                        name='iframeEditor'
                    />
                </div>
            )}
        </>
    );
}
