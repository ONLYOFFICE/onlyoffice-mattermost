// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createStore} from 'redux';

import {OPEN_EDITOR_MODAL, UPDATE_HEALTH_STATUS} from '../actions/types';

import rootReducer from './index';

describe('root reducer', () => {
    it('combines plugin slices', () => {
        const store = createStore(rootReducer);
        const state = store.getState();

        expect(state).toEqual(expect.objectContaining({
            permissionsModal: expect.any(Object),
            editorModal: expect.any(Object),
            managerModal: expect.any(Object),
            converterModal: expect.any(Object),
            health: expect.any(Object),
        }));

        store.dispatch({type: OPEN_EDITOR_MODAL, payload: {id: '1'}});
        store.dispatch({type: UPDATE_HEALTH_STATUS, payload: false});

        expect(store.getState().editorModal.isVisible).toBe(true);
        expect(store.getState().health.healthy).toBe(false);
    });
});
