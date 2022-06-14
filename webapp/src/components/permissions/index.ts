/**
 *
 * (c) Copyright Ascensio System SIA 2022
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
