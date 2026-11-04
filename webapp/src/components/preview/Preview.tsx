// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2026
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

import {openEditorInNewWindow} from 'util/editor';
import fileHelper from 'util/file';
import {getTranslations} from 'util/lang';

import manifest from 'manifest';
import editor from 'public/images/editor.svg';
import editorDark from 'public/images/editor_dark.svg';
import permissions from 'public/images/permissions.svg';
import permissionsDark from 'public/images/permissions_dark.svg';
import React from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {openPermissions} from 'redux/actions';

import type {FileInfo} from '@mattermost/types/files';

import FileDetailsRows from 'components/preview/FileDetails';

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
    const isHealthy = useSelector((state: any) => {
        const pluginState = state['plugins-' + manifest.id] || {};
        const health = pluginState.health || {healthy: true};
        const healthy = health.healthy !== false;
        return healthy;
    });

    const showPermissions = fileHelper.isExtensionSupported(props.fileInfo.extension, true) && fileHelper.isFileAuthor(props.fileInfo) && isHealthy;

    return (
        <div
            className='file-details__container'
            data-theme={props.theme}
            data-dark-theme={props.darkTheme}
        >
            <div className='file-details__preview onlyoffice_preview__left'>
                <img
                    className='onlyoffice_preview__icon'
                    alt='file preview'
                    src={icon}
                />
                <div
                    className='file-details__onlyoffice onlyoffice_preview__actions'
                    data-theme={props.theme}
                    data-dark-theme={props.darkTheme}
                >
                    {isHealthy && (
                        <button
                            type='button'
                            className='onlyoffice_preview__btn onlyoffice_preview__btn_open'
                            title={i18n['preview.open_button']}
                            onClick={() => openEditorInNewWindow(props.fileInfo, props.theme)}
                            data-theme={props.theme}
                        >
                            <img
                                alt=''
                                src={props.theme === 'dark' ? editorDark : editor}
                            />
                            <span>{i18n['preview.open_button']}</span>
                        </button>
                    )}
                    {showPermissions && (
                        <button
                            type='button'
                            className='onlyoffice_preview__btn onlyoffice_preview__btn_share'
                            title={i18n['plugin.access_button']}
                            onClick={() => openPermissions(props.fileInfo)(dispatch)}
                            data-theme={props.theme}
                        >
                            <img
                                alt=''
                                src={props.theme === 'dark' ? permissionsDark : permissions}
                            />
                            <span>{i18n['preview.share_button']}</span>
                        </button>
                    )}
                    {!isHealthy && (
                        <div
                            className='onlyoffice_preview__disabled'
                            title={(i18n as any)['preview.server_unavailable'] || 'ONLYOFFICE document server is currently unavailable'}
                            style={{
                                opacity: 0.5,
                                cursor: 'not-allowed',
                                fontSize: '12px',
                                color: props.theme === 'dark' ? '#999' : '#666',
                                marginTop: '4px',
                            }}
                        >
                            {(i18n as any)['preview.server_unavailable'] || 'Server unavailable'}
                        </div>
                    )}
                </div>
            </div>
            <div
                className='file-details'
                style={{display: 'flex', flexDirection: 'column'}}
            >
                <div className='file-details__name'>{props.fileInfo.name}</div>
                <div className='file-details__info'>
                    {`${i18n['preview.file_type']} ${props.fileInfo.extension.toUpperCase()}`}
                </div>
                <FileDetailsRows
                    fileId={props.fileInfo.id}
                    theme={props.theme}
                />
            </div>
        </div>
    );
}
