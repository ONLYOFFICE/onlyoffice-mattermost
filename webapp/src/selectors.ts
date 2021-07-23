import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from './manifest';

//@ts-ignore
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const isEditorModalVisible = (state: GlobalState) => getPluginState(state).editorModalVisible.isVisible;
export const editorModalFileInfo = (state: GlobalState) => getPluginState(state).editorModalVisible.fileInfo;

export const isPermissionsModalVisible = (state: GlobalState) => getPluginState(state).permissionsModalVisible.isVisible;
export const permissionsModalFileInfo = (state: GlobalState) => getPluginState(state).permissionsModalVisible.fileInfo;
