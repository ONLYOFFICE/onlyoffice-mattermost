// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {
    CLOSE_CONVERTER_MODAL,
    CLOSE_EDITOR_MODAL,
    CLOSE_MANAGER_MODAL,
    CLOSE_PERMISSIONS_MODAL,
    OPEN_CONVERTER_MODAL,
    OPEN_EDITOR_MODAL,
    OPEN_MANAGER_MODAL,
    OPEN_PERMISSIONS_MODAL,
    UPDATE_HEALTH_STATUS,
} from './types';

import {
    closeConverter,
    closeEditor,
    closeManager,
    closePermissions,
    openConverter,
    openEditor,
    openManager,
    openPermissions,
    updateHealthStatus,
} from './index';

describe('redux actions', () => {
    const fileInfo = {id: 'file-1'} as any;
    const dispatch = jest.fn();

    beforeEach(() => {
        dispatch.mockClear();
    });

    it.each([
        [openEditor, OPEN_EDITOR_MODAL, fileInfo],
        [openConverter, OPEN_CONVERTER_MODAL, fileInfo],
        [openPermissions, OPEN_PERMISSIONS_MODAL, fileInfo],
    ])('%p dispatches open with payload', (action, type, payload) => {
        action(payload)(dispatch);
        expect(dispatch).toHaveBeenCalledWith({type, payload});
    });

    it.each([
        [closeEditor, CLOSE_EDITOR_MODAL],
        [closeConverter, CLOSE_CONVERTER_MODAL],
        [closePermissions, CLOSE_PERMISSIONS_MODAL],
        [closeManager, CLOSE_MANAGER_MODAL],
        [openManager, OPEN_MANAGER_MODAL],
    ])('%p dispatches expected type', (action, type) => {
        action()(dispatch);
        expect(dispatch).toHaveBeenCalledWith({type});
    });

    it('updateHealthStatus dispatches payload', () => {
        updateHealthStatus(false)(dispatch);
        expect(dispatch).toHaveBeenCalledWith({
            type: UPDATE_HEALTH_STATUS,
            payload: false,
        });
    });
});
