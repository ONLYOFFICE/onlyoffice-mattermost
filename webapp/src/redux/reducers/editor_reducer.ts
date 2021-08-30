/**
 *
 * (c) Copyright Ascensio System SIA 2021
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
