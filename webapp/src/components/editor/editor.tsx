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

import React, {useCallback} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import {Dispatch} from 'redux';

import {ONLYOFFICE_PLUGIN_API, ONLYOFFICE_PLUGIN_API_EDITOR} from 'utils';

import {EditorLoader} from './editor_loader';

interface EditorProps {
    visible: boolean;
    close: () => (dispatch: Dispatch) => void;
    fileInfo?: FileInfo;
}

const Editor = ({visible, close, fileInfo}: EditorProps) => {
    const lang = localStorage.getItem('onlyoffice_locale') || 'en';
    const handleClose = useCallback(() => {
        if (!visible) {
            return;
        }
        const editorBackdrop = document.getElementById('editor-backdrop');

        if (editorBackdrop) {
            editorBackdrop.classList.add('onlyoffice-modal__backdrop_hide');
        }

        setTimeout(() => close(), 300);
    }, [close, visible]);

    const onEscape = useCallback(
        (event) => {
            if (event.keyCode === 27) {
                handleClose();
            }
        },
        [handleClose],
    );

    React.useEffect(() => {
        if (!visible || !fileInfo) {
            return;
        }
        window.addEventListener('ONLYOFFICE_CLOSED', handleClose);
        document.addEventListener('keydown', onEscape, false);

        // eslint-disable-next-line consistent-return
        return () => {
            window.removeEventListener('ONLYOFFICE_CLOSED', handleClose);
            document.removeEventListener('keydown', onEscape, false);
        };
    }, [fileInfo, visible, handleClose, onEscape]);

    return (
        <>
            {visible && (
                <div
                    id='editor-backdrop'
                    className='onlyoffice-modal__backdrop'
                >
                    <EditorLoader/>
                    <iframe
                        src={`${ONLYOFFICE_PLUGIN_API}${ONLYOFFICE_PLUGIN_API_EDITOR}?file=${fileInfo?.id}&lang=${lang}`}
                        className='onlyoffice-modal__frame'
                        name='iframeEditor'
                    />
                </div>
            )}
        </>
    );
};

export default Editor;
