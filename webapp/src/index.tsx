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

import manifest from 'manifest';
import React from 'react';
import Reducer from 'redux/reducers';
import {openConverter, openEditor, openManager, openPermissions} from 'redux/actions';
import type {Action, AnyAction, Store} from 'redux';
import type {FileInfo} from 'mattermost-redux/types/files';
import type {GlobalState} from 'mattermost-redux/types/store';
import type {ThunkDispatch} from 'redux-thunk';

import {getTranslations} from 'util/lang';
import {isConvertSupported, isExtensionSupported, isFileAuthor} from 'util/file';

import {ManagerIcon} from 'components/manager/Icon';
import OnlyofficeEditor from 'components/editor';
import OnlyofficeFileConverter from 'components/converter';
import OnlyofficeFilePermissions from 'components/permissions';
import OnlyofficeFilePreview from 'components/preview';
import OnlyofficeManager from 'components/manager';

import 'public/scss/icons.scss';
import 'public/scss/editor.scss';

export default class Plugin {
    public async initialize(registry: any, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerTranslations(getTranslations);
        registry.registerReducer(Reducer);
        registry.registerRootComponent(OnlyofficeEditor);
        registry.registerRootComponent(OnlyofficeFilePermissions);
        registry.registerRootComponent(OnlyofficeManager);
        registry.registerRootComponent(OnlyofficeFileConverter);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;

        if (registry.registerFileDropdownMenuAction) {
            registry.registerFileDropdownMenuAction(
                (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension),
                () => getTranslations()['plugin.open_button'],
                (fileInfo: FileInfo) => dispatch(openEditor(fileInfo)),
            );
            registry.registerFileDropdownMenuAction(
                (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension, true) && isFileAuthor(fileInfo),
                () => getTranslations()['plugin.access_button'],
                (fileInfo: FileInfo) => dispatch(openPermissions(fileInfo)),
            );
            registry.registerFileDropdownMenuAction(
                (fileInfo: FileInfo) => isConvertSupported(fileInfo.extension) && isFileAuthor(fileInfo),
                () => getTranslations()['plugin.convert_button'],
                (fileInfo: FileInfo) => dispatch(openConverter(fileInfo)),
            );
        }

        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                return isExtensionSupported(fileInfo.extension) && fileInfo.extension !== 'pdf';
            },
            OnlyofficeFilePreview,
        );

        registry.registerFileUploadMethod(
            <ManagerIcon/>,
            () => dispatch(openManager()),
            'ONLYOFFICE',
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
