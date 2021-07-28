/* eslint-disable no-console */
import {AnyAction, Store} from 'redux';
import {GlobalState} from 'mattermost-redux/types/store';
import {FileInfo} from 'mattermost-redux/types/files';
import {ThunkDispatch} from 'redux-thunk';

import {openEditor, openPermissions} from 'actions';

import {isExtensionSupported, isFileAuthor} from 'utils/file_utils';

import manifest from './manifest';
import Reducer from './reducer';
import Editor from './components/editor';
import Permissions from './components/permissions';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public async initialize(registry: any, store: Store<GlobalState>) {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(Editor);
        registry.registerRootComponent(Permissions);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;
        registry.registerFileDropdownMenuAction(
            (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension),
            'Open file in ONLYOFFICE',
            (fileInfo: FileInfo) => dispatch(openEditor(fileInfo)),
        );
        registry.registerFileDropdownMenuAction(
            (fileInfo: FileInfo) => isExtensionSupported(fileInfo.extension) && isFileAuthor(fileInfo),
            'Change access rights',
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
