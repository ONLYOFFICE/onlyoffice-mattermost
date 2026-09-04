// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {id as pluginName} from 'manifest';

import * as httpModule from './http';

import {get, getHealthStatus, getPluginConfig, post} from './index';

describe('api helpers', () => {
    afterEach(() => {
        jest.restoreAllMocks();
    });

    it('get and post delegate to http with methods', async () => {
        const spy = jest.spyOn(httpModule, 'http').mockResolvedValue({ok: true});

        await get('/path', {credentials: 'include'});
        expect(spy).toHaveBeenCalledWith('/path', expect.objectContaining({
            method: 'GET',
            credentials: 'include',
        }));

        await post('/path', {a: 1});
        expect(spy).toHaveBeenCalledWith('/path', expect.objectContaining({
            method: 'POST',
            body: JSON.stringify({a: 1}),
        }));
    });

    it('loads plugin config and health status', async () => {
        const spy = jest.spyOn(httpModule, 'http').mockResolvedValue({healthy: true});

        await getPluginConfig();
        expect(spy).toHaveBeenCalledWith(`/plugins/${pluginName}/api/config`, expect.any(Object));

        await getHealthStatus();
        expect(spy).toHaveBeenCalledWith(`/plugins/${pluginName}/api/health`, expect.any(Object));
    });
});
