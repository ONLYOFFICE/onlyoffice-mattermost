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

import {isExtensionSupported, isFileAuthor} from 'util/file';
import {getTranslations} from 'util/lang';

import manifest from 'manifest';
import type {Action, AnyAction, Store} from 'redux';
import {openEditor, openPermissions} from 'redux/actions';
import Reducer from 'redux/reducers';
import type {ThunkDispatch} from 'redux-thunk';

import type {FileInfo} from 'mattermost-redux/types/files';
import type {GlobalState} from 'mattermost-redux/types/store';

import OnlyofficeEditor from 'components/editor';
import OnlyofficeFilePermissions from 'components/permissions';
import OnlyofficeFilePreview from 'components/preview';

import 'public/scss/icons.scss';
import 'public/scss/editor.scss';

export default class Plugin {
    public async initialize(registry: any, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerTranslations(getTranslations);
        registry.registerReducer(Reducer);
        registry.registerRootComponent(OnlyofficeEditor);
        registry.registerRootComponent(OnlyofficeFilePermissions);
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
        }

        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                return isExtensionSupported(fileInfo.extension) && fileInfo.extension !== 'pdf';
            },
            OnlyofficeFilePreview,
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
