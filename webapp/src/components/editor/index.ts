import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {isEditorModalVisible, editorModalFileInfo} from 'redux/selectors';

import {closeEditor} from 'redux/actions';

import Editor from './editor';

const mapStateToProps = (state: GlobalState) => ({
    visible: isEditorModalVisible(state),
    fileInfo: editorModalFileInfo(state),
});

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    close: closeEditor,
}, dispatch);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export default connect(mapStateToProps, mapDispatchToProps)(Editor as any);
