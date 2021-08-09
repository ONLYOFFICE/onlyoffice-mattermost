import {Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {CLOSE_EDITOR_MODAL, OPEN_EDITOR_MODAL} from './action_types';

export const closeEditor = () => (dispatch: Dispatch) => {
    dispatch({
        type: CLOSE_EDITOR_MODAL,
    });
};

export const openEditor = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: OPEN_EDITOR_MODAL,
        payload: fileInfo,
    });
};
