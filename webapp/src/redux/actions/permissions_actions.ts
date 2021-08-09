import {FileInfo} from 'mattermost-redux/types/files';
import {Dispatch} from 'redux';

import {CLOSE_PERMISSIONS_MODAL, OPEN_PERMISSIONS_MODAL} from './action_types';

export const closePermissions = () => (dispatch: Dispatch) => {
    dispatch({
        type: CLOSE_PERMISSIONS_MODAL,
    });
};

export const openPermissions = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: OPEN_PERMISSIONS_MODAL,
        payload: fileInfo,
    });
};
