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

import fileHelper from 'util/file';
import {getTranslations} from 'util/lang';

import editor from 'public/images/editor.svg';
import editorDark from 'public/images/editor_dark.svg';
import permissions from 'public/images/permissions.svg';
import permissionsDark from 'public/images/permissions_dark.svg';
import React from 'react';
import {useDispatch} from 'react-redux';
import {openEditor, openPermissions} from 'redux/actions';

import type {FileInfo} from 'mattermost-redux/types/files';

import 'public/scss/preview.scss';

type Props = {
    fileInfo: FileInfo;
    theme: string;
    darkTheme: string | undefined;
}

export default function OnlyofficeFilePreview(props: Props) {
    const i18n = getTranslations();
    const dispatch = useDispatch();
    const icon = fileHelper.getIconByExt(props.fileInfo.extension);
    const showPermissions = fileHelper.isExtensionSupported(props.fileInfo.extension, true) && fileHelper.isFileAuthor(props.fileInfo);

    return (
        <div
            className='file-details__container'
            data-theme={props.theme}
            data-dark-theme={props.darkTheme}
        >
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
                <div
                    className='file-details__onlyoffice'
                    data-theme={props.theme}
                    data-dark-theme={props.darkTheme}
                >
                    {
                        showPermissions &&
                            (
                                <img
                                    className='onlyoffice_preview__btn'
                                    alt={'permissions button'}
                                    onClick={() => openPermissions(props.fileInfo)(dispatch)}
                                    src={props.theme === 'dark' ? permissionsDark : permissions}
                                    data-theme={props.theme}
                                    data-dark-theme={props.darkTheme}
                                />
                            )
                    }
                    <img
                        className='onlyoffice_preview__btn'
                        alt={'open editor'}
                        onClick={() => openEditor(props.fileInfo)(dispatch)}
                        src={props.theme === 'dark' ? editorDark : editor}
                        data-theme={props.theme}
                        data-dark-theme={props.darkTheme}
                    />
                </div>
            </div>
        </div>
    );
}
