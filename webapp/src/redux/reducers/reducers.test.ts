// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {converterModal} from './converter';
import {editorModal} from './editor';
import health from './health';
import {managerModal} from './manager';
import {permissionsModal} from './permissions';

import {
    CLOSE_CONVERTER_MODAL,
    CLOSE_EDITOR_MODAL,
    CLOSE_MANAGER_MODAL,
    CLOSE_PERMISSIONS_MODAL,
    MATTERMOST_ME_ACTION,
    OPEN_CONVERTER_MODAL,
    OPEN_EDITOR_MODAL,
    OPEN_MANAGER_MODAL,
    OPEN_PERMISSIONS_MODAL,
    UPDATE_HEALTH_STATUS,
} from '../actions/types';

describe('redux reducers', () => {
    const fileInfo = {id: 'file-1', name: 'doc.docx'};

    it('toggles editor modal', () => {
        expect(editorModal(undefined, {type: 'unknown'})).toEqual({isVisible: false, fileInfo: null});
        expect(editorModal(undefined, {type: OPEN_EDITOR_MODAL, payload: fileInfo})).toEqual({
            isVisible: true,
            fileInfo,
        });
        expect(editorModal({isVisible: true, fileInfo: fileInfo as any}, {type: CLOSE_EDITOR_MODAL})).toEqual({
            isVisible: false,
            fileInfo: null,
        });
    });

    it('toggles converter modal', () => {
        expect(converterModal(undefined, {type: OPEN_CONVERTER_MODAL, payload: fileInfo}).isVisible).toBe(true);
        expect(converterModal({isVisible: true, fileInfo: fileInfo as any}, {type: CLOSE_CONVERTER_MODAL}).isVisible).toBe(false);
    });

    it('toggles permissions modal and handles me action', () => {
        expect(permissionsModal(undefined, {type: OPEN_PERMISSIONS_MODAL, payload: fileInfo})).toEqual({
            isVisible: true,
            fileInfo,
        });
        expect(permissionsModal({isVisible: true, fileInfo: fileInfo as any}, {type: CLOSE_PERMISSIONS_MODAL})).toEqual({
            isVisible: false,
            fileInfo: null,
        });

        const state = {isVisible: false, fileInfo: null};
        expect(permissionsModal(state, {type: MATTERMOST_ME_ACTION, data: {locale: 'ru'}})).toBe(state);
        expect(window.localStorage.getItem('onlyoffice_locale')).toBe('ru');
    });

    it('toggles manager modal', () => {
        expect(managerModal(undefined, {type: OPEN_MANAGER_MODAL})).toEqual({isVisible: true});
        expect(managerModal({isVisible: true}, {type: CLOSE_MANAGER_MODAL})).toEqual({isVisible: false});
    });

    it('updates health status', () => {
        expect(health(undefined, {type: 'unknown'})).toEqual({healthy: true});
        expect(health(undefined, {type: UPDATE_HEALTH_STATUS, payload: false})).toEqual({healthy: false});
        expect(health({healthy: false}, {type: UPDATE_HEALTH_STATUS, payload: true})).toEqual({healthy: true});
    });
});
