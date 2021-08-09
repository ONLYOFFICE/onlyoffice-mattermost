import {AnyAction} from 'redux';

import {OPEN_PERMISSIONS_MODAL, CLOSE_PERMISSIONS_MODAL} from 'redux/actions/action_types';

export const permissionsModalVisible = (state = {isVisible: false, fileInfo: null}, action: AnyAction) => {
    switch (action.type) {
    case OPEN_PERMISSIONS_MODAL:
        return {
            isVisible: true,
            fileInfo: action.payload,
        };
    case CLOSE_PERMISSIONS_MODAL:
        return {
            isVisible: false,
            fileInfo: null,
        };
    default:
        return state;
    }
};
