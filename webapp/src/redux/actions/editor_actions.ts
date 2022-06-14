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

import {Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {CLOSE_EDITOR_MODAL, OPEN_EDITOR_MODAL} from './action_types';

export const closeEditor = () => (dispatch: Dispatch) => {
    dispatch({
        type: CLOSE_EDITOR_MODAL,
    });
};

export const openEditor = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: OPEN_EDITOR_MODAL,
        payload: fileInfo,
    });
};
