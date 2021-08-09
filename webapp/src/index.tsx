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
import 'public/scss/classes/modal_editor.scss';

export default class Plugin {
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

        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                return isExtensionSupported(fileInfo.extension) && fileInfo.extension !== 'pdf';
            },
            FilePreviewOverride,
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
