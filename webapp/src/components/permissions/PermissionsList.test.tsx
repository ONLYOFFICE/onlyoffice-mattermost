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

import {FileAccess} from 'util/permission';

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import {PermissionsList} from './PermissionsList';

const invokeStyles = (styles: Record<string, any> = {}) => {
    Object.values(styles).forEach((styleFn) => {
        if (typeof styleFn === 'function') {
            styleFn({}, {isFocused: false});
            styleFn({}, {isFocused: true});
        }
    });
};

jest.mock('react-select', () => ({
    __esModule: true,
    default: ({styles, onChange, options}: any) => {
        invokeStyles(styles);
        return (
            <button
                type='button'
                onClick={() => onChange?.(options?.[0])}
            >
                {'perm-select'}
            </button>
        );
    },
}));

describe('PermissionsList', () => {
    const users = [{
        value: 'u1',
        label: 'alice',
        avatarUrl: 'avatar.png',
        fileAccess: FileAccess.EDIT_ONLY,
        email: 'a@ex.com',
    }];

    it('renders fetch error state', () => {
        render(
            <PermissionsList
                error={true}
                users={[]}
                theme='light'
                darkTheme={undefined}
                onRemoveUser={jest.fn()}
                onChangeUserPermissions={jest.fn()}
            />,
        );

        expect(screen.getByText(/could not fetch users/i)).toBeInTheDocument();
    });

    it('renders users and removes on action', async () => {
        const onRemoveUser = jest.fn();
        render(
            <PermissionsList
                error={false}
                users={users}
                theme='light'
                darkTheme={undefined}
                onRemoveUser={onRemoveUser}
                onChangeUserPermissions={jest.fn()}
            />,
        );

        expect(screen.getByText('@alice')).toBeInTheDocument();

        const removeButtons = screen.getAllByRole('button', {name: /close/i});
        await userEvent.click(removeButtons[removeButtons.length - 1]);
        expect(onRemoveUser).toHaveBeenCalledWith('alice');
    });

    it.each([
        ['dark', 'indigo'],
        ['dark', 'onyx'],
        ['dark', undefined],
    ] as const)('covers %s/%s permission select styles and change', async (theme, darkTheme) => {
        const onChangeUserPermissions = jest.fn();
        render(
            <PermissionsList
                error={false}
                users={users}
                theme={theme}
                darkTheme={darkTheme}
                onRemoveUser={jest.fn()}
                onChangeUserPermissions={onChangeUserPermissions}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: 'perm-select'}));
        expect(onChangeUserPermissions).toHaveBeenCalledWith('alice', expect.any(String));
    });
});
