// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {getCookie} from './cookie';

describe('getCookie', () => {
    beforeEach(() => {
        document.cookie.split(';').forEach((cookie) => {
            const name = cookie.split('=')[0]?.trim();
            if (name) {
                document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
            }
        });
    });

    it('returns empty string when cookie is missing', () => {
        expect(getCookie('MMUSERID')).toBe('');
    });

    it('returns cookie value when present', () => {
        document.cookie = 'MMUSERID=user-1';
        expect(getCookie('MMUSERID')).toBe('user-1');
    });
});
