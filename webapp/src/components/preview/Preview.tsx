/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
import {useDispatch} from 'react-redux';
import {FileInfo} from 'mattermost-redux/types/files';

import {openEditor, openPermissions} from 'redux/actions';

import fileHelper from 'util/file';
import {getTranslations} from 'util/lang';

import editor from 'public/images/editor.svg';
import permissions from 'public/images/permissions.svg';
import 'public/scss/preview.scss';

type Props = {
    fileInfo: FileInfo;
}

export default function OnlyofficeFilePreview(props: Props) {
    const i18n = getTranslations();
    const dispatch = useDispatch();
    const icon = fileHelper.getIconByExt(props.fileInfo.extension);
    const showPermissions = fileHelper.isExtensionSupported(props.fileInfo.extension, true) && fileHelper.isFileAuthor(props.fileInfo);

    return (
        <div className='file-details__container'>
            <a
                className='file-details__preview'
                onClick={(e) => e.preventDefault()}
            >
                <a className='file-details__preview-helper'/>
                <img
                    alt='file preview'
                    src={icon}
                />
            </a>
            <div
                className='file-details'
                style={{display: 'flex', flexDirection: 'column'}}
            >
                <div className='file-details__name'>{props.fileInfo.name}</div>
                <div className='file-details__info'>
                    {`${i18n['preview.file_type']} ${props.fileInfo.extension.toUpperCase()}`}
                </div>
                <div className='file-details__onlyoffice'>
                    {
                        showPermissions &&
                            (
                                <img
                                    className='onlyoffice_preview__btn'
                                    alt={'permissions button'}
                                    onClick={() => dispatch(openPermissions(props.fileInfo))}
                                    src={permissions}
                                />
                            )
                    }
                    <img
                        className='onlyoffice_preview__btn'
                        alt={'open editor'}
                        onClick={() => dispatch(openEditor(props.fileInfo))}
                        src={editor}
                    />
                </div>
            </div>
        </div>
    );
}
