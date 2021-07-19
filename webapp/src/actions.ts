import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {id as pluginId} from './manifest';
import {OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL} from './action_types';

export const closeRootModal = () => (dispatch: Dispatch) => {
    dispatch({
        type: CLOSE_ROOT_MODAL,
    });
};

export const postDropdownMenuAction = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: OPEN_ROOT_MODAL,
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
