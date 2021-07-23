import {AnyAction, combineReducers} from 'redux';

import {OPEN_EDITOR_MODAL, CLOSE_EDITOR_MODAL, OPEN_PERMISSIONS_MODAL, CLOSE_PERMISSIONS_MODAL} from './action_types';

const editorModalVisible = (state = {isVisible: false, fileInfo: null}, action: AnyAction) => {
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

const permissionsModalVisible = (state = {isVisible: false, fileInfo: null}, action: AnyAction) => {
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

export default combineReducers({
    editorModalVisible,
    permissionsModalVisible,
});

