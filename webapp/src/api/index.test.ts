// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
