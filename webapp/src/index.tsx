import { AnyAction, Store } from 'redux';

import { GlobalState } from 'mattermost-redux/types/store';
import { FileInfo } from 'mattermost-redux/types/files';
import { ThunkDispatch } from 'redux-thunk';

import { postDropdownMenuAction } from 'actions';

import manifest from './manifest';
import Reducer from './reducer';
import Root from './root';

// eslint-disable-next-line import/no-unresolved

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: any, store: Store<GlobalState>) {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(Root);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;
        registry.registerFileDropdownMenuAction(
            () => true,
            "ONLYOFFICE",
            (fileInfo: FileInfo) => dispatch(postDropdownMenuAction(fileInfo))
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
