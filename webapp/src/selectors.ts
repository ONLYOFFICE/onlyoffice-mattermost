import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from './manifest';

//@ts-ignore
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};
export const isRootModalVisible = (state: GlobalState) => getPluginState(state).rootModalVisible.isVisible;
export const fileInfo = (state: GlobalState) => getPluginState(state).rootModalVisible.fileInfo;
