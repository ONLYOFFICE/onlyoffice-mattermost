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
import {AnyAction} from 'redux';

import {OPEN_PERMISSIONS_MODAL, CLOSE_PERMISSIONS_MODAL, MATTERMOST_ME_ACTION} from 'redux/actions/types';
import {getTranslations} from 'util/lang';

export const permissionsModal = (state = {isVisible: false, fileInfo: null}, action: AnyAction) => {
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
    case MATTERMOST_ME_ACTION:
        getTranslations(action.data.locale || 'en');
        return state;
    default:
        return state;
    }
};
