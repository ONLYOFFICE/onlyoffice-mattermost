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

import {Client4} from 'mattermost-redux/client';

import {debounce, debounceUsersLoad, pipe} from './func';

jest.mock('mattermost-redux/client', () => ({
    Client4: {
        searchUsers: jest.fn(),
    },
}));

describe('func utils', () => {
    beforeEach(() => {
        jest.useFakeTimers();
    });

    afterEach(() => {
        jest.useRealTimers();
    });

    it('debounces callback invocation', () => {
        const cb = jest.fn();
        const debounced = debounce(cb, 100);
        debounced('a');
        debounced('b');

        expect(cb).not.toHaveBeenCalled();

        jest.advanceTimersByTime(100);

        expect(cb).toHaveBeenCalledTimes(1);
        expect(cb).toHaveBeenCalledWith('b');
    });

    it('pipes functions left to right', () => {
        const double = (n: number) => n * 2;
        const inc = (n: number) => n + 1;
        expect(pipe(double, inc)(3)).toBe(7);
    });

    it('debounceUsersLoad searches channel users', async () => {
        (Client4.searchUsers as jest.Mock).mockResolvedValue([
            {id: 'author', username: 'author'},
            {id: 'u2', username: 'bob', email: 'b@ex.com'},
        ]);
        const callback = jest.fn();
        const load = debounceUsersLoad(
            {id: 'ch1', team_id: 't1'} as any,
            {user_id: 'author'} as any,
            [],
        );

        load('', callback);

        jest.advanceTimersByTime(2000);
        expect(Client4.searchUsers).not.toHaveBeenCalled();

        load('bo', callback);

        jest.advanceTimersByTime(2000);

        await Promise.resolve();
        await Promise.resolve();

        expect(Client4.searchUsers).toHaveBeenCalled();
        expect(callback).toHaveBeenCalledWith([
            expect.objectContaining({value: 'u2', label: 'bob'}),
        ]);
    });
});
