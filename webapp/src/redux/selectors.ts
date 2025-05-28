// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

import type {GlobalState} from 'mattermost-redux/types/store';
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import {id as pluginId} from '../manifest';

//@ts-expect-error: Suppressing error because state['plugins-' + pluginId] might be undefined
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const editorModalVisible = (state: GlobalState) => getPluginState(state).editorModal.isVisible;
export const editorModalFileInfo = (state: GlobalState) => getPluginState(state).editorModal.fileInfo;

export const permissionsModalVisible = (state: GlobalState) => getPluginState(state).permissionsModal.isVisible;
export const permissionsModalFileInfo = (state: GlobalState) => getPluginState(state).permissionsModal.fileInfo;

export const managerModalVisible = (state: GlobalState) => getPluginState(state).managerModal.isVisible;

export const getCurrentTheme = (state: GlobalState) => {
    const theme = getTheme(state);
    const dark = theme.type === 'indigo' || theme.type === 'onyx';
    return dark ? 'dark' : 'light';
};
