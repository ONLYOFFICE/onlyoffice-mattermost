import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {isPermissionsModalVisible, permissionsModalFileInfo} from 'redux/selectors';

import {closePermissions} from 'redux/actions';

import Permissions from './permissions';

const mapStateToProps = (state: GlobalState) => ({
    visible: isPermissionsModalVisible(state),
    fileInfo: permissionsModalFileInfo(state),
});

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    close: closePermissions,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Permissions);
