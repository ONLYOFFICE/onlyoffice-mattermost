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
import {AnyAction, Store} from 'redux';
import {GlobalState} from 'mattermost-redux/types/store';
import {FileInfo} from 'mattermost-redux/types/files';
import {ThunkDispatch} from 'redux-thunk';

import {isExtensionSupported, isFileAuthor} from 'utils/file';

import FilePreviewOverride from 'components/file_preview/file_preview';

import {openEditor, openPermissions} from 'redux/actions';

import Reducer from 'redux/reducers';

import manifest from 'manifest';

import Editor from 'components/editor';
import Permissions from 'components/permissions';
import 'public/scss/icons.scss';
import 'public/scss/modal_editor.scss';
import {getTranslations} from 'utils/i18n';
import {apiHealth} from 'api';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public async initialize(registry: any, store: Store<GlobalState>) {
        //TODO: Plugin lifetime healthchecks
        const isAlive = await apiHealth();
        if (!isAlive) {
            return;
        }
        registry.registerTranslations(getTranslations);
        const i18n = getTranslations();
        registry.registerReducer(Reducer);
        registry.registerRootComponent(Editor);
        registry.registerRootComponent(Permissions);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;
        registry.registerFileDropdownMenuAction(
            (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension),
            i18n['plugin.open_button'],
            (fileInfo: FileInfo) => dispatch(openEditor(fileInfo)),
        );

        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                return isExtensionSupported(fileInfo.extension) && fileInfo.extension !== 'pdf';
            },
            FilePreviewOverride,
        );

        registry.registerFileDropdownMenuAction(
            (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension, true) && isFileAuthor(fileInfo),
            i18n['plugin.access_button'],
            (fileInfo: FileInfo) => dispatch(openPermissions(fileInfo)),
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
