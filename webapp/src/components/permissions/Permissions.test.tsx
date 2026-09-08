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

import {ONLYOFFICE_WILDCARD_USER} from 'util/const';

import {render, screen, waitFor} from '@testing-library/react';
import * as api from 'api';
import React from 'react';

import Permissions from './Permissions';

jest.mock('api', () => ({
    get: jest.fn(),
    ONLYOFFICE_PLUGIN_PERMISSIONS: '/api/permissions',
}));

jest.mock('mattermost-redux/client', () => ({
    Client4: {
        getTeamByName: jest.fn(),
        getChannelByName: jest.fn(),
    },
}));

jest.mock('components/permissions/PermissionsHeader', () => ({
    PermissionsHeader: () => <div data-testid='permissions-header'/>,
}));

jest.mock('components/permissions/PermissionsList', () => ({
    PermissionsList: () => <div data-testid='permissions-list'/>,
}));

jest.mock('components/permissions/PermissionsFooter', () => ({
    PermissionsFooter: () => <div data-testid='permissions-footer'/>,
}));

describe('Permissions', () => {
    const close = jest.fn(() => jest.fn());
    const fileInfo = {id: 'file-1'} as any;

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders nothing when not visible', () => {
        const {container} = render(
            <Permissions
                visible={false}
                close={close}
                fileInfo={fileInfo}
                editorOpen={false}
                theme='light'
                darkTheme={undefined}
            />,
        );

        expect(container).toBeEmptyDOMElement();
    });

    it('loads permissions and renders modal sections', async () => {
        (api.get as jest.Mock).mockResolvedValue([
            {id: ONLYOFFICE_WILDCARD_USER, username: '*', email: '*', permissions: {edit: false}},
            {id: 'u1', username: 'alice', email: 'a@ex.com', permissions: {edit: true}},
        ]);

        render(
            <Permissions
                visible={true}
                close={close}
                fileInfo={fileInfo}
                editorOpen={false}
                theme='light'
                darkTheme={undefined}
            />,
        );

        expect(await screen.findByTestId('permissions-header')).toBeInTheDocument();
        expect(screen.getByTestId('permissions-footer')).toBeInTheDocument();

        await waitFor(() => expect(api.get).toHaveBeenCalled());
    });

    it('shows footer in error path when fetch fails', async () => {
        (api.get as jest.Mock).mockRejectedValue(new Error('network'));

        render(
            <Permissions
                visible={true}
                close={close}
                fileInfo={fileInfo}
                editorOpen={false}
                theme='light'
                darkTheme={undefined}
            />,
        );

        expect(await screen.findByTestId('permissions-footer')).toBeInTheDocument();
    });
});
