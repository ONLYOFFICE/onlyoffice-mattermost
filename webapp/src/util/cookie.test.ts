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
