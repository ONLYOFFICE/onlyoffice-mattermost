import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {closeRootModal} from 'actions';
import {isRootModalVisible, fileInfo} from 'selectors';

import Root from './root';

const mapStateToProps = (state: {visible: boolean, fileInfo: FileInfo}) => ({
    visible: isRootModalVisible(state),
    fileInfo: fileInfo(state),
});

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    close: closeRootModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Root);
