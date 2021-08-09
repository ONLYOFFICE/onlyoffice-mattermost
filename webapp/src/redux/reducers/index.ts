import {combineReducers} from 'redux';

import {editorModalVisible} from './editor_reducer';
import {permissionsModalVisible} from './permissions_reducer';

export default combineReducers({
    editorModalVisible,
    permissionsModalVisible,
});

