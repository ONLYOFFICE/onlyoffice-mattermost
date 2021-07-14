import {combineReducers} from 'redux';

import {OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL} from './action_types';

const rootModalVisible = (state = {isVisible: false, fileInfo: null}, action: any) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return {
            isVisible: true,
            fileInfo: action.payload
        };
    case CLOSE_ROOT_MODAL:
        return {
            isVisible: false,
            fileInfo: null
        };
    default:
        return state;
    }
};

export default combineReducers({
    rootModalVisible
});

