// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {FileAccess} from 'util/permission';

import {render, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import {PermissionsHeader} from './PermissionsHeader';

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
    default: ({styles, onChange, options, isDisabled}: any) => {
        invokeStyles(styles);
        return (
            <button
                type='button'
                disabled={isDisabled}
                onClick={() => onChange?.(options?.[0])}
            >
                {'wildcard-select'}
            </button>
        );
    },
}));

jest.mock('react-select/async', () => ({
    __esModule: true,
    default: ({styles, onChange, isDisabled, loadOptions}: any) => {
        invokeStyles(styles);
        loadOptions?.('alice', jest.fn());
        return (
            <button
                type='button'
                disabled={isDisabled}
                onClick={() => onChange?.([{
                    value: 'u2',
                    label: 'bob',
                    avatarUrl: '',
                    fileAccess: 'Edit',
                    email: 'b@ex.com',
                }])}
            >
                {'user-select'}
            </button>
        );
    },
}));

jest.mock('util/func', () => ({
    debounceUsersLoad: () => jest.fn(),
}));

describe('PermissionsHeader', () => {
    const baseProps = {
        loading: false,
        channel: {id: 'ch1', team_id: 't1'} as any,
        fileInfo: {id: 'f1', user_id: 'author'} as any,
        wildcardAccess: FileAccess.READ_ONLY,
        users: [],
        theme: 'light',
        darkTheme: undefined as string | undefined,
        onSetWildcardAccess: jest.fn(),
        onAppendUsers: jest.fn(),
    };

    beforeEach(() => {
        jest.clearAllMocks();
        window.history.pushState({}, '', '/team/channels/town-square');
    });

    it('adds selected users and updates wildcard access', async () => {
        render(<PermissionsHeader {...baseProps}/>);

        await userEvent.click(screen.getByRole('button', {name: 'user-select'}));
        await userEvent.click(screen.getByRole('button', {name: /add/i}));

        expect(baseProps.onAppendUsers).toHaveBeenCalledWith([
            expect.objectContaining({value: 'u2', label: 'bob'}),
        ]);

        await userEvent.click(screen.getByRole('button', {name: 'wildcard-select'}));
        expect(baseProps.onSetWildcardAccess).toHaveBeenCalled();
    });

    it('renders compact mode without channel user picker', async () => {
        window.history.pushState({}, '', '/team/messages/@alice');
        render(
            <PermissionsHeader
                {...baseProps}
                channel={null}
            />,
        );

        expect(screen.queryByRole('button', {name: 'user-select'})).not.toBeInTheDocument();

        await waitFor(() => {
            expect(screen.getByText(/access rights/i)).toBeInTheDocument();
        });
    });

    it.each([
        ['dark', 'indigo'],
        ['dark', 'onyx'],
        ['dark', undefined],
    ] as const)('applies %s/%s theme styles', (theme, darkTheme) => {
        render(
            <PermissionsHeader
                {...baseProps}
                theme={theme}
                darkTheme={darkTheme}
            />,
        );

        expect(screen.getByRole('button', {name: 'wildcard-select'})).toBeInTheDocument();
    });
});
