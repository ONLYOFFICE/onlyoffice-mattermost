/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

import {FileInfo} from 'mattermost-redux/types/files';
import {useDispatch} from 'react-redux';

import {getIconByExt, isExtensionSupported, isFileAuthor} from 'utils/file';
import {getTranslations} from 'utils/i18n';

import {openEditor, openPermissions} from 'redux/actions';

import 'public/scss/preview.scss';
import editor from 'public/images/editor.svg';
import permissions from 'public/images/permissions.svg';

type Props = {
    fileInfo: FileInfo;
};

export default function FilePreviewOverride(props: Props) {
    const dispatch = useDispatch();
    const icon = getIconByExt(props.fileInfo.extension);
    const permissionsWindow = isExtensionSupported(props.fileInfo.extension, true) && isFileAuthor(props.fileInfo);

    //A temporary solution to close mm-6.0 modal
    const modalCloseBtn = document.querySelector('div.file-preview-modal-main-actions__actions > button') as HTMLButtonElement;
    return (
        <div className='modal-image-backround'>
            <div className='modal-image__content'>
                <div className='file-details__container'>
                    <a
                        className='file-details__preview'
                        onClick={(e) => e.preventDefault()}
                    >
                        <span className='file-details__preview-helper'/>
                        <img
                            alt='file preview'
                            src={icon}
                        />
                    </a>
                    <div
                        className='file-details'
                        style={{position: 'relative'}}
                    >
                        <div className='file-details__name'>{props.fileInfo.name}</div>
                        <div className='file-details__info'>{`File type ${props.fileInfo.extension.toUpperCase()}`}</div>
                        {
                            permissionsWindow &&
                            (
                                <img
                                    className='onlyoffice_preview__btn'
                                    style={{position: 'absolute', right: '6.5rem', bottom: '2rem'}}
                                    alt={getTranslations()['plugin.access_button']}
                                    onClick={() => dispatch(openPermissions(props.fileInfo))}
                                    src={permissions}
                                />
                            )
                        }
                        <img
                            className='onlyoffice_preview__btn'
                            style={{position: 'absolute', right: '2rem', bottom: '2rem'}}
                            alt={getTranslations()['preview.open_button']}
                            onClick={() => {
                                if (modalCloseBtn) {
                                    modalCloseBtn.click();
                                }
                                dispatch(openEditor(props.fileInfo));
                            }}
                            src={editor}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
