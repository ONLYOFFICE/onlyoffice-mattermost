import {AnyAction} from 'redux';

import {OPEN_EDITOR_MODAL, CLOSE_EDITOR_MODAL} from 'redux/actions/action_types';

export const editorModalVisible = (state = {isVisible: false, fileInfo: null}, action: AnyAction) => {
    switch (action.type) {
    case OPEN_EDITOR_MODAL:
        return {
            isVisible: true,
            fileInfo: action.payload,
        };
    case CLOSE_EDITOR_MODAL:
        return {
            isVisible: false,
            fileInfo: null,
        };
    default:
        return state;
    }
};
