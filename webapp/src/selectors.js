import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible.isVisible;
export const fileInfo = (state) => getPluginState(state).rootModalVisible.fileInfo;
