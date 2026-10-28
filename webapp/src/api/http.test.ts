// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
