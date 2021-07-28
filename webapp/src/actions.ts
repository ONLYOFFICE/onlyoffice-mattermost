import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {id as pluginId} from './manifest';

import {OPEN_EDITOR_MODAL, CLOSE_EDITOR_MODAL,
    CLOSE_PERMISSIONS_MODAL, OPEN_PERMISSIONS_MODAL} from './action_types';

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

export const closePermissions = () => (dispatch: Dispatch) => {
    dispatch({
        type: CLOSE_PERMISSIONS_MODAL,
    });
};

export const openPermissions = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: OPEN_PERMISSIONS_MODAL,
        payload: fileInfo,
    });
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const getPluginServerRoute = (state: any) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};
