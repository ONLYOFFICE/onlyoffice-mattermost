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

import {render, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as api from 'api';
import React from 'react';

import Converter from './Converter';

jest.mock('api', () => ({
    get: jest.fn(),
    post: jest.fn(),
    ONLYOFFICE_PLUGIN_GET_CODE: '/api/code',
    ONLYOFFICE_PLUGIN_CONVERT: '/api/convert',
}));

describe('Converter', () => {
    const close = jest.fn(() => jest.fn());
    const fileInfo = {id: 'file-1', name: 'doc.doc'} as any;

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders nothing when hidden', () => {
        const {container} = render(
            <Converter
                visible={false}
                fileInfo={fileInfo}
                theme='light'
                darkTheme={undefined}
                close={close}
            />,
        );

        expect(container).toBeEmptyDOMElement();
    });

    it('closes after successful convert', async () => {
        jest.useFakeTimers();
        const user = userEvent.setup({advanceTimers: jest.advanceTimersByTime});
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue({error: 0});

        render(
            <Converter
                visible={true}
                fileInfo={fileInfo}
                theme='light'
                darkTheme={undefined}
                close={close}
            />,
        );

        await user.click(screen.getByRole('button', {name: /convert/i}));
        await waitFor(() => expect(api.post).toHaveBeenCalled());

        jest.advanceTimersByTime(300);

        await waitFor(() => expect(close).toHaveBeenCalled());

        jest.useRealTimers();
    });

    it('asks for password when convert returns -5', async () => {
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue({error: -5});

        render(
            <Converter
                visible={true}
                fileInfo={fileInfo}
                theme='light'
                darkTheme={undefined}
                close={close}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /convert/i}));
        expect(await screen.findByPlaceholderText(/password/i)).toBeInTheDocument();
    });

    it('asks for format when convert returns -9', async () => {
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue({error: -9});

        render(
            <Converter
                visible={true}
                fileInfo={fileInfo}
                theme='light'
                darkTheme={undefined}
                close={close}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /convert/i}));
        expect(await screen.findByText(/select output format/i)).toBeInTheDocument();
    });

    it('shows convert failure for other errors', async () => {
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue({error: -1});

        render(
            <Converter
                visible={true}
                fileInfo={fileInfo}
                theme='light'
                darkTheme={undefined}
                close={close}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /convert/i}));
        expect(await screen.findByText(/failed to convert/i)).toBeInTheDocument();
    });
});
