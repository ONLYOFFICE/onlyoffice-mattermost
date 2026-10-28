// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {FileAccess} from 'util/permission';

import {render, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as api from 'api';
import React from 'react';

import {PermissionsFooter} from './PermissionsFooter';

jest.mock('api', () => ({
    get: jest.fn(),
    post: jest.fn(),
    ONLYOFFICE_PLUGIN_GET_CODE: '/api/code',
    ONLYOFFICE_PLUGIN_PERMISSIONS: '/api/permissions',
}));

describe('PermissionsFooter', () => {
    const onClose = jest.fn();
    const fileInfo = {id: 'file-1'} as any;
    const users = [{
        value: 'u1',
        label: 'alice',
        avatarUrl: '',
        fileAccess: FileAccess.EDIT_ONLY,
        email: 'a@ex.com',
    }];

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('closes without posting when cancel is clicked', async () => {
        render(
            <PermissionsFooter
                fileInfo={fileInfo}
                loading={false}
                users={users}
                wildcardAccess={FileAccess.READ_ONLY}
                onClose={onClose}
                theme='light'
                darkTheme={undefined}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /cancel/i}));
        expect(onClose).toHaveBeenCalled();
        expect(api.post).not.toHaveBeenCalled();
    });

    it('posts permissions and closes on save', async () => {
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue(undefined);

        render(
            <PermissionsFooter
                fileInfo={fileInfo}
                loading={false}
                users={users}
                wildcardAccess={FileAccess.READ_ONLY}
                onClose={onClose}
                theme='light'
                darkTheme={undefined}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /save/i}));
        await waitFor(() => expect(api.post).toHaveBeenCalled());
        expect(onClose).toHaveBeenCalled();
    });
});
