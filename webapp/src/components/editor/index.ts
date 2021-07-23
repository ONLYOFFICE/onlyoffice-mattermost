import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {closeEditor} from 'actions';
import {isEditorModalVisible, editorModalFileInfo} from 'selectors';

import Editor from './editor';

const mapStateToProps = (state: GlobalState) => ({
    visible: isEditorModalVisible(state),
    fileInfo: editorModalFileInfo(state),
});

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    close: closeEditor,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Editor);
