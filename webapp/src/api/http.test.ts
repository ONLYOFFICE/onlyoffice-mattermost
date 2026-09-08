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

import {http} from './http';

describe('http', () => {
    const originalFetch = global.fetch;

    afterEach(() => {
        global.fetch = originalFetch;
    });

    it('returns parsed JSON on success', async () => {
        global.fetch = jest.fn().mockResolvedValue({
            ok: true,
            json: async () => ({ok: true}),
        }) as any;

        await expect(http('/api', {method: 'GET'})).resolves.toEqual({ok: true});
    });

    it('throws on non-ok responses', async () => {
        global.fetch = jest.fn().mockResolvedValue({
            ok: false,
            statusText: 'Forbidden',
            json: async () => ({}),
        }) as any;

        await expect(http('/api', {method: 'GET'})).rejects.toThrow('Forbidden');
    });

    it('returns undefined when body is not JSON', async () => {
        global.fetch = jest.fn().mockResolvedValue({
            ok: true,
            json: async () => {
                throw new Error('invalid json');
            },
        }) as any;

        await expect(http('/api', {method: 'GET'})).resolves.toBeUndefined();
    });
});
