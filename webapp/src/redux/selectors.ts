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

import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from '../manifest';

//@ts-ignore
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const isEditorModalVisible = (state: GlobalState) => getPluginState(state).editorModalVisible.isVisible;
export const editorModalFileInfo = (state: GlobalState) => getPluginState(state).editorModalVisible.fileInfo;

export const isPermissionsModalVisible = (state: GlobalState) => getPluginState(state).permissionsModalVisible.isVisible;
export const permissionsModalFileInfo = (state: GlobalState) => getPluginState(state).permissionsModalVisible.fileInfo;
